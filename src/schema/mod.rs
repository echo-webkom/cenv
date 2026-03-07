use regex::Regex;
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::net::IpAddr;

#[derive(Debug, Serialize, Deserialize)]
pub struct Schema {
    pub entries: Vec<Entry>,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Entry {
    /// The entry key, all caps.
    pub key: String,
    /// Description of this entry. Describes what the value should be in a human readable way.
    pub hint: Option<String>,
    /// If true, field must be present and have a non-empty value.
    pub required: bool,
    /// Optional default value for field to be filled in when generating env.
    pub default: Option<String>,
    /// Optional list of legal field values.
    /// If Some then the fields value must be one of the given.
    pub legal_values: Option<Vec<String>>,
    /// Optional required length of field.
    pub required_length: Option<usize>,
    /// Optional regex string the fields value must match.
    pub regex_match: Option<String>,
    /// Optional type specifier for field.
    pub kind: Option<EntryKind>,
}

#[derive(Debug, Serialize, Deserialize)]
pub enum EntryKind {
    Integer { min: Option<i64>, max: Option<i64> },
    Float { min: Option<f64>, max: Option<f64> },
    String,
    Url,
    Email,
    Bool,
    IpAddress,
    Path,
}

/// Generates .env file content based on the given schema and existing env values.
///
/// For keys with existing non-empty values, those values are preserved.
/// For keys without values, the schema's default value is used if available,
/// otherwise the value is left empty.
pub fn generate_env(schema: &Schema, existing_env: &HashMap<String, String>) -> String {
    let mut output = String::new();

    for entry in &schema.entries {
        // Determine the value to use:
        // 1. Use existing value from .env if present and non-empty
        // 2. Otherwise use default from schema if present
        // 3. Otherwise leave empty
        let value = match existing_env.get(&entry.key) {
            Some(existing_value) if !existing_value.is_empty() => existing_value.clone(),
            _ => entry.default.clone().unwrap_or_default(),
        };

        output.push_str(&format!("{}={}\n", entry.key, value));
    }

    output
}

/// A validation error for a specific entry.
#[derive(Debug)]
pub struct ValidationError {
    pub key: String,
    pub message: String,
    pub hint: Option<String>,
}

impl std::fmt::Display for ValidationError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match &self.hint {
            Some(hint) => write!(f, "{}: {}\n\thint: {}", self.key, self.message, hint),
            None => write!(f, "{}: {}", self.key, self.message),
        }
    }
}

/// Validates environment variables against a schema.
pub fn validate_env(schema: &Schema, env_vars: &HashMap<String, String>) -> Vec<ValidationError> {
    let mut errors = Vec::new();

    for entry in &schema.entries {
        let value = env_vars.get(&entry.key);
        let hint = entry.hint.clone();

        // Check required
        if entry.required {
            match value {
                None => {
                    errors.push(ValidationError {
                        key: entry.key.clone(),
                        message: "required field is missing".to_string(),
                        hint,
                    });
                    continue;
                }
                Some(v) if v.is_empty() => {
                    errors.push(ValidationError {
                        key: entry.key.clone(),
                        message: "required field is empty".to_string(),
                        hint,
                    });
                    continue;
                }
                _ => {}
            }
        }

        // Skip further validation if value is missing or empty (and not required)
        let value = match value {
            Some(v) if !v.is_empty() => v,
            _ => continue,
        };

        // Check required_length
        if let Some(required_length) = entry.required_length
            && value.len() != required_length
        {
            errors.push(ValidationError {
                key: entry.key.clone(),
                message: format!("expected length {}, got {}", required_length, value.len()),
                hint: hint.clone(),
            });
        }

        // Check legal_values
        if let Some(ref legal_values) = entry.legal_values
            && !legal_values.contains(&value.to_string())
        {
            errors.push(ValidationError {
                key: entry.key.clone(),
                message: format!(
                    "value '{}' is not one of: {}",
                    value,
                    legal_values.join(", ")
                ),
                hint: hint.clone(),
            });
        }

        // Check regex_match
        if let Some(ref regex_pattern) = entry.regex_match {
            match Regex::new(regex_pattern) {
                Ok(re) => {
                    if !re.is_match(value) {
                        errors.push(ValidationError {
                            key: entry.key.clone(),
                            message: format!("value does not match pattern: {}", regex_pattern),
                            hint: hint.clone(),
                        });
                    }
                }
                Err(e) => {
                    errors.push(ValidationError {
                        key: entry.key.clone(),
                        message: format!("invalid regex pattern: {}", e),
                        hint: hint.clone(),
                    });
                }
            }
        }

        // Check kind-specific validation
        if let Some(ref kind) = entry.kind
            && let Some(err) = validate_kind(value, kind)
        {
            errors.push(ValidationError {
                key: entry.key.clone(),
                message: err,
                hint: hint.clone(),
            });
        }
    }

    errors
}

/// Validates a value against its EntryKind.
fn validate_kind(value: &str, kind: &EntryKind) -> Option<String> {
    match kind {
        EntryKind::Integer { min, max } => match value.parse::<i64>() {
            Ok(n) => {
                if let Some(min_val) = min
                    && n < *min_val
                {
                    return Some(format!("value {} is less than minimum {}", n, min_val));
                }
                if let Some(max_val) = max
                    && n > *max_val
                {
                    return Some(format!("value {} is greater than maximum {}", n, max_val));
                }
                None
            }
            Err(_) => Some(format!("'{}' is not a valid integer", value)),
        },
        EntryKind::Float { min, max } => match value.parse::<f64>() {
            Ok(n) => {
                if let Some(min_val) = min
                    && n < *min_val
                {
                    return Some(format!("value {} is less than minimum {}", n, min_val));
                }
                if let Some(max_val) = max
                    && n > *max_val
                {
                    return Some(format!("value {} is greater than maximum {}", n, max_val));
                }
                None
            }
            Err(_) => Some(format!("'{}' is not a valid float", value)),
        },
        EntryKind::String => None,
        EntryKind::Url => {
            // URL validation using RFC 3986 compliant regex
            // Matches: scheme://[userinfo@]host[:port][/path][?query][#fragment]
            // Also supports file:// URLs with absolute paths
            let url_regex = Regex::new(
                r"(?i)^[a-z][a-z0-9+.-]*://(([a-z0-9._~%!$&'()*+,;=:-]*@)?([a-z0-9.-]+|\[[a-f0-9:]+\])(:[0-9]+)?)?(/[a-z0-9._~%!$&'()*+,;=:@/-]*)?(\?[a-z0-9._~%!$&'()*+,;=:@/?-]*)?(\#[a-z0-9._~%!$&'()*+,;=:@/?-]*)?$"
            ).unwrap();
            if url_regex.is_match(value) {
                None
            } else {
                Some(format!("'{}' is not a valid URL", value))
            }
        }
        EntryKind::Email => {
            // Email validation using RFC 5321/5322 compliant regex
            // Allows: local-part@domain where domain has at least one dot
            let email_regex = Regex::new(
                r"(?i)^[a-z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-z0-9]([a-z0-9-]*[a-z0-9])?(\.[a-z0-9]([a-z0-9-]*[a-z0-9])?)+$"
            ).unwrap();
            if email_regex.is_match(value) {
                None
            } else {
                Some(format!("'{}' is not a valid email address", value))
            }
        }
        EntryKind::Bool => match value.to_lowercase().as_str() {
            "true" | "false" | "1" | "0" | "yes" | "no" => None,
            _ => Some(format!("'{}' is not a valid boolean", value)),
        },
        EntryKind::IpAddress => match value.parse::<IpAddr>() {
            Ok(_) => None,
            Err(_) => Some(format!("'{}' is not a valid IP address", value)),
        },
        EntryKind::Path => {
            // Path validation - validates Unix and Windows path formats
            // Unix: starts with / or ./ or ../ or is relative
            // Windows: starts with drive letter (C:\) or UNC (\\server\share) or relative
            let unix_path_regex =
                Regex::new(r"^(/|\.{1,2}/)?([a-zA-Z0-9._-]+/?)*[a-zA-Z0-9._-]*$").unwrap();
            let windows_path_regex = Regex::new(
                r"(?i)^([a-z]:[/\\]|\\\\[a-z0-9._-]+[/\\][a-z0-9._-]+)?([a-zA-Z0-9._-]+[/\\]?)*[a-zA-Z0-9._-]*$"
            ).unwrap();

            if value.is_empty() {
                Some("path cannot be empty".to_string())
            } else if unix_path_regex.is_match(value) || windows_path_regex.is_match(value) {
                None
            } else {
                Some(format!("'{}' is not a valid path", value))
            }
        }
    }
}

#[cfg(test)]
mod tests;
