use super::*;

// ==================== Helper Functions ====================

fn create_entry(key: &str, entry: Entry) -> Entry {
    Entry {
        key: key.to_string(),
        ..entry
    }
}

fn create_schema(entries: Vec<Entry>) -> Schema {
    Schema { entries }
}

fn default_entry() -> Entry {
    Entry {
        key: String::new(),
        hint: None,
        required: false,
        default: None,
        legal_values: None,
        required_length: None,
        regex_match: None,
        kind: None,
    }
}

// Alias for backwards compatibility with tests
fn default_options() -> Entry {
    default_entry()
}

// ==================== generate_env Tests ====================

mod generate_env_tests {
    use super::*;

    #[test]
    fn empty_schema_produces_empty_output() {
        let schema = create_schema(vec![]);
        let existing = HashMap::new();
        let output = generate_env(&schema, &existing);
        assert_eq!(output, "");
    }

    #[test]
    fn single_entry_no_existing_no_default() {
        let schema = create_schema(vec![create_entry("FOO", default_options())]);
        let existing = HashMap::new();
        let output = generate_env(&schema, &existing);
        assert_eq!(output, "FOO=\n");
    }

    #[test]
    fn single_entry_with_default_no_existing() {
        let mut opts = default_options();
        opts.default = Some("bar".to_string());
        let schema = create_schema(vec![create_entry("FOO", opts)]);
        let existing = HashMap::new();
        let output = generate_env(&schema, &existing);
        assert_eq!(output, "FOO=bar\n");
    }

    #[test]
    fn existing_value_takes_precedence_over_default() {
        let mut opts = default_options();
        opts.default = Some("default_value".to_string());
        let schema = create_schema(vec![create_entry("FOO", opts)]);
        let mut existing = HashMap::new();
        existing.insert("FOO".to_string(), "existing_value".to_string());
        let output = generate_env(&schema, &existing);
        assert_eq!(output, "FOO=existing_value\n");
    }

    #[test]
    fn empty_existing_value_uses_default() {
        let mut opts = default_options();
        opts.default = Some("default_value".to_string());
        let schema = create_schema(vec![create_entry("FOO", opts)]);
        let mut existing = HashMap::new();
        existing.insert("FOO".to_string(), "".to_string());
        let output = generate_env(&schema, &existing);
        assert_eq!(output, "FOO=default_value\n");
    }

    #[test]
    fn empty_existing_value_no_default_stays_empty() {
        let schema = create_schema(vec![create_entry("FOO", default_options())]);
        let mut existing = HashMap::new();
        existing.insert("FOO".to_string(), "".to_string());
        let output = generate_env(&schema, &existing);
        assert_eq!(output, "FOO=\n");
    }

    #[test]
    fn multiple_entries_preserve_order() {
        let schema = create_schema(vec![
            create_entry("AAA", default_options()),
            create_entry("BBB", default_options()),
            create_entry("CCC", default_options()),
        ]);
        let existing = HashMap::new();
        let output = generate_env(&schema, &existing);
        assert_eq!(output, "AAA=\nBBB=\nCCC=\n");
    }

    #[test]
    fn existing_values_for_some_entries() {
        let mut opts_with_default = default_options();
        opts_with_default.default = Some("default_b".to_string());

        let schema = create_schema(vec![
            create_entry("AAA", default_options()),
            create_entry("BBB", opts_with_default),
            create_entry("CCC", default_options()),
        ]);
        let mut existing = HashMap::new();
        existing.insert("AAA".to_string(), "value_a".to_string());
        existing.insert("CCC".to_string(), "value_c".to_string());
        let output = generate_env(&schema, &existing);
        assert_eq!(output, "AAA=value_a\nBBB=default_b\nCCC=value_c\n");
    }

    #[test]
    fn extra_existing_values_are_ignored() {
        let schema = create_schema(vec![create_entry("FOO", default_options())]);
        let mut existing = HashMap::new();
        existing.insert("FOO".to_string(), "foo_value".to_string());
        existing.insert("BAR".to_string(), "bar_value".to_string());
        existing.insert("BAZ".to_string(), "baz_value".to_string());
        let output = generate_env(&schema, &existing);
        assert_eq!(output, "FOO=foo_value\n");
    }

    #[test]
    fn values_with_special_characters() {
        let schema = create_schema(vec![create_entry("FOO", default_options())]);
        let mut existing = HashMap::new();
        existing.insert(
            "FOO".to_string(),
            "value with spaces & special=chars!".to_string(),
        );
        let output = generate_env(&schema, &existing);
        assert_eq!(output, "FOO=value with spaces & special=chars!\n");
    }
}

// ==================== validate_env Tests ====================

mod validate_env_tests {
    use super::*;

    // -------------------- Hint Tests --------------------

    mod hint_tests {
        use super::*;

        #[test]
        fn error_includes_hint_when_present() {
            let mut opts = default_options();
            opts.required = true;
            opts.hint = Some("Should be your API key from the dashboard".to_string());
            let schema = create_schema(vec![create_entry("API_KEY", opts)]);
            let env = HashMap::new();
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert_eq!(
                errors[0].hint,
                Some("Should be your API key from the dashboard".to_string())
            );
            // Check Display includes hint
            let display = format!("{}", errors[0]);
            assert!(display.contains("hint:"));
            assert!(display.contains("Should be your API key from the dashboard"));
        }

        #[test]
        fn error_without_hint_displays_without_hint() {
            let mut opts = default_options();
            opts.required = true;
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let env = HashMap::new();
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert_eq!(errors[0].hint, None);
            // Check Display does not include hint
            let display = format!("{}", errors[0]);
            assert!(!display.contains("hint:"));
        }

        #[test]
        fn hint_included_in_kind_validation_error() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Integer {
                min: None,
                max: None,
            });
            opts.hint = Some("Must be a valid port number".to_string());
            let schema = create_schema(vec![create_entry("PORT", opts)]);
            let mut env = HashMap::new();
            env.insert("PORT".to_string(), "not_a_number".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert_eq!(
                errors[0].hint,
                Some("Must be a valid port number".to_string())
            );
        }

        #[test]
        fn hint_included_in_regex_validation_error() {
            let mut opts = default_options();
            opts.regex_match = Some(r"^\d{3}-\d{4}$".to_string());
            opts.hint = Some("Format: XXX-XXXX".to_string());
            let schema = create_schema(vec![create_entry("PHONE", opts)]);
            let mut env = HashMap::new();
            env.insert("PHONE".to_string(), "invalid".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert_eq!(errors[0].hint, Some("Format: XXX-XXXX".to_string()));
        }
    }

    // -------------------- Required Field Tests --------------------

    mod required_tests {
        use super::*;

        #[test]
        fn required_field_present_and_non_empty_passes() {
            let mut opts = default_options();
            opts.required = true;
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let mut env = HashMap::new();
            env.insert("FOO".to_string(), "value".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn required_field_missing_fails() {
            let mut opts = default_options();
            opts.required = true;
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let env = HashMap::new();
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert_eq!(errors[0].key, "FOO");
            assert!(errors[0].message.contains("missing"));
        }

        #[test]
        fn required_field_empty_fails() {
            let mut opts = default_options();
            opts.required = true;
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let mut env = HashMap::new();
            env.insert("FOO".to_string(), "".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert_eq!(errors[0].key, "FOO");
            assert!(errors[0].message.contains("empty"));
        }

        #[test]
        fn non_required_field_missing_passes() {
            let schema = create_schema(vec![create_entry("FOO", default_options())]);
            let env = HashMap::new();
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn non_required_field_empty_passes() {
            let schema = create_schema(vec![create_entry("FOO", default_options())]);
            let mut env = HashMap::new();
            env.insert("FOO".to_string(), "".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }
    }

    // -------------------- Required Length Tests --------------------

    mod required_length_tests {
        use super::*;

        #[test]
        fn exact_length_passes() {
            let mut opts = default_options();
            opts.required_length = Some(5);
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let mut env = HashMap::new();
            env.insert("FOO".to_string(), "12345".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn too_short_fails() {
            let mut opts = default_options();
            opts.required_length = Some(5);
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let mut env = HashMap::new();
            env.insert("FOO".to_string(), "1234".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert!(errors[0].message.contains("expected length 5"));
            assert!(errors[0].message.contains("got 4"));
        }

        #[test]
        fn too_long_fails() {
            let mut opts = default_options();
            opts.required_length = Some(5);
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let mut env = HashMap::new();
            env.insert("FOO".to_string(), "123456".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert!(errors[0].message.contains("expected length 5"));
            assert!(errors[0].message.contains("got 6"));
        }

        #[test]
        fn zero_length_required() {
            let mut opts = default_options();
            opts.required_length = Some(0);
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let mut env = HashMap::new();
            env.insert("FOO".to_string(), "a".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
        }

        #[test]
        fn empty_value_skips_length_check() {
            let mut opts = default_options();
            opts.required_length = Some(5);
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let mut env = HashMap::new();
            env.insert("FOO".to_string(), "".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn missing_value_skips_length_check() {
            let mut opts = default_options();
            opts.required_length = Some(5);
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let env = HashMap::new();
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn unicode_characters_counted_correctly() {
            let mut opts = default_options();
            opts.required_length = Some(3);
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let mut env = HashMap::new();
            // "héllo" has 5 characters but 6 bytes
            env.insert("FOO".to_string(), "héy".to_string());
            let errors = validate_env(&schema, &env);
            // Rust's len() counts bytes, not characters
            // "héy" is 4 bytes (h=1, é=2, y=1)
            assert_eq!(errors.len(), 1);
        }
    }

    // -------------------- Legal Values Tests --------------------

    mod legal_values_tests {
        use super::*;

        #[test]
        fn value_in_legal_values_passes() {
            let mut opts = default_options();
            opts.legal_values = Some(vec!["a".to_string(), "b".to_string(), "c".to_string()]);
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let mut env = HashMap::new();
            env.insert("FOO".to_string(), "b".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn value_not_in_legal_values_fails() {
            let mut opts = default_options();
            opts.legal_values = Some(vec!["a".to_string(), "b".to_string(), "c".to_string()]);
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let mut env = HashMap::new();
            env.insert("FOO".to_string(), "d".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert!(errors[0].message.contains("not one of"));
        }

        #[test]
        fn case_sensitive_legal_values() {
            let mut opts = default_options();
            opts.legal_values = Some(vec!["Foo".to_string()]);
            let schema = create_schema(vec![create_entry("VAR", opts)]);
            let mut env = HashMap::new();
            env.insert("VAR".to_string(), "foo".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
        }

        #[test]
        fn empty_legal_values_list() {
            let mut opts = default_options();
            opts.legal_values = Some(vec![]);
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let mut env = HashMap::new();
            env.insert("FOO".to_string(), "anything".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
        }

        #[test]
        fn single_legal_value() {
            let mut opts = default_options();
            opts.legal_values = Some(vec!["only".to_string()]);
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let mut env = HashMap::new();
            env.insert("FOO".to_string(), "only".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn empty_value_skips_legal_values_check() {
            let mut opts = default_options();
            opts.legal_values = Some(vec!["a".to_string()]);
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let mut env = HashMap::new();
            env.insert("FOO".to_string(), "".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }
    }

    // -------------------- Regex Match Tests --------------------

    mod regex_match_tests {
        use super::*;

        #[test]
        fn matching_regex_passes() {
            let mut opts = default_options();
            opts.regex_match = Some(r"^\d{3}-\d{4}$".to_string());
            let schema = create_schema(vec![create_entry("PHONE", opts)]);
            let mut env = HashMap::new();
            env.insert("PHONE".to_string(), "123-4567".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn non_matching_regex_fails() {
            let mut opts = default_options();
            opts.regex_match = Some(r"^\d{3}-\d{4}$".to_string());
            let schema = create_schema(vec![create_entry("PHONE", opts)]);
            let mut env = HashMap::new();
            env.insert("PHONE".to_string(), "12-34567".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert!(errors[0].message.contains("does not match pattern"));
        }

        #[test]
        fn invalid_regex_pattern_fails() {
            let mut opts = default_options();
            opts.regex_match = Some(r"[invalid".to_string());
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let mut env = HashMap::new();
            env.insert("FOO".to_string(), "anything".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert!(errors[0].message.contains("invalid regex"));
        }

        #[test]
        fn uuid_v4_regex() {
            let mut opts = default_options();
            opts.regex_match = Some(
                r"^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
                    .to_string(),
            );
            let schema = create_schema(vec![create_entry("UUID", opts)]);

            let mut env = HashMap::new();
            env.insert(
                "UUID".to_string(),
                "550e8400-e29b-41d4-a716-446655440000".to_string(),
            );
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn empty_value_skips_regex_check() {
            let mut opts = default_options();
            opts.regex_match = Some(r"^\d+$".to_string());
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let mut env = HashMap::new();
            env.insert("FOO".to_string(), "".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn anchored_vs_unanchored_regex() {
            let mut opts = default_options();
            opts.regex_match = Some(r"\d+".to_string()); // Not anchored
            let schema = create_schema(vec![create_entry("FOO", opts)]);
            let mut env = HashMap::new();
            env.insert("FOO".to_string(), "abc123def".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty()); // Should pass because regex is not anchored
        }
    }

    // -------------------- Integer Kind Tests --------------------

    mod integer_kind_tests {
        use super::*;

        #[test]
        fn valid_integer_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Integer {
                min: None,
                max: None,
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "42".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn negative_integer_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Integer {
                min: None,
                max: None,
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "-42".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn invalid_integer_fails() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Integer {
                min: None,
                max: None,
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "not_a_number".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert!(errors[0].message.contains("not a valid integer"));
        }

        #[test]
        fn float_as_integer_fails() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Integer {
                min: None,
                max: None,
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "3.14".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
        }

        #[test]
        fn integer_within_min_max_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Integer {
                min: Some(1),
                max: Some(100),
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "50".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn integer_at_min_boundary_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Integer {
                min: Some(1),
                max: Some(100),
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "1".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn integer_at_max_boundary_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Integer {
                min: Some(1),
                max: Some(100),
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "100".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn integer_below_min_fails() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Integer {
                min: Some(1),
                max: Some(100),
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "0".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert!(errors[0].message.contains("less than minimum"));
        }

        #[test]
        fn integer_above_max_fails() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Integer {
                min: Some(1),
                max: Some(100),
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "101".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert!(errors[0].message.contains("greater than maximum"));
        }

        #[test]
        fn integer_with_only_min() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Integer {
                min: Some(0),
                max: None,
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "999999".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn integer_with_only_max() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Integer {
                min: None,
                max: Some(100),
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "-999999".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn integer_with_negative_bounds() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Integer {
                min: Some(-100),
                max: Some(-1),
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "-50".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }
    }

    // -------------------- Float Kind Tests --------------------

    mod float_kind_tests {
        use super::*;

        #[test]
        fn valid_float_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Float {
                min: None,
                max: None,
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "3.14".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn integer_as_float_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Float {
                min: None,
                max: None,
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "42".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn negative_float_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Float {
                min: None,
                max: None,
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "-3.14".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn scientific_notation_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Float {
                min: None,
                max: None,
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "1.5e10".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn invalid_float_fails() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Float {
                min: None,
                max: None,
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "not_a_number".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert!(errors[0].message.contains("not a valid float"));
        }

        #[test]
        fn float_within_bounds_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Float {
                min: Some(0.0),
                max: Some(1.0),
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "0.5".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn float_below_min_fails() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Float {
                min: Some(0.0),
                max: Some(1.0),
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "-0.1".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert!(errors[0].message.contains("less than minimum"));
        }

        #[test]
        fn float_above_max_fails() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Float {
                min: Some(0.0),
                max: Some(1.0),
            });
            let schema = create_schema(vec![create_entry("NUM", opts)]);
            let mut env = HashMap::new();
            env.insert("NUM".to_string(), "1.1".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert!(errors[0].message.contains("greater than maximum"));
        }
    }

    // -------------------- String Kind Tests --------------------

    mod string_kind_tests {
        use super::*;

        #[test]
        fn any_string_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::String);
            let schema = create_schema(vec![create_entry("STR", opts)]);
            let mut env = HashMap::new();
            env.insert("STR".to_string(), "anything at all!".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn special_characters_pass() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::String);
            let schema = create_schema(vec![create_entry("STR", opts)]);
            let mut env = HashMap::new();
            env.insert(
                "STR".to_string(),
                "!@#$%^&*()_+-=[]{}|;':\",./<>?".to_string(),
            );
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }
    }

    // -------------------- URL Kind Tests --------------------

    mod url_kind_tests {
        use super::*;

        #[test]
        fn http_url_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Url);
            let schema = create_schema(vec![create_entry("URL", opts)]);
            let mut env = HashMap::new();
            env.insert("URL".to_string(), "http://example.com".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn https_url_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Url);
            let schema = create_schema(vec![create_entry("URL", opts)]);
            let mut env = HashMap::new();
            env.insert(
                "URL".to_string(),
                "https://example.com/path?query=1".to_string(),
            );
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn postgres_url_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Url);
            let schema = create_schema(vec![create_entry("URL", opts)]);
            let mut env = HashMap::new();
            env.insert(
                "URL".to_string(),
                "postgres://user:pass@localhost:5432/db".to_string(),
            );
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn missing_scheme_fails() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Url);
            let schema = create_schema(vec![create_entry("URL", opts)]);
            let mut env = HashMap::new();
            env.insert("URL".to_string(), "example.com".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert!(errors[0].message.contains("not a valid URL"));
        }

        #[test]
        fn file_url_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Url);
            let schema = create_schema(vec![create_entry("URL", opts)]);
            let mut env = HashMap::new();
            env.insert("URL".to_string(), "file:///path/to/file".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }
    }

    // -------------------- Email Kind Tests --------------------

    mod email_kind_tests {
        use super::*;

        #[test]
        fn valid_email_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Email);
            let schema = create_schema(vec![create_entry("EMAIL", opts)]);
            let mut env = HashMap::new();
            env.insert("EMAIL".to_string(), "user@example.com".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn email_with_subdomain_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Email);
            let schema = create_schema(vec![create_entry("EMAIL", opts)]);
            let mut env = HashMap::new();
            env.insert("EMAIL".to_string(), "user@mail.example.com".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn email_with_plus_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Email);
            let schema = create_schema(vec![create_entry("EMAIL", opts)]);
            let mut env = HashMap::new();
            env.insert("EMAIL".to_string(), "user+tag@example.com".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn missing_at_sign_fails() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Email);
            let schema = create_schema(vec![create_entry("EMAIL", opts)]);
            let mut env = HashMap::new();
            env.insert("EMAIL".to_string(), "userexample.com".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert!(errors[0].message.contains("not a valid email"));
        }

        #[test]
        fn missing_domain_dot_fails() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Email);
            let schema = create_schema(vec![create_entry("EMAIL", opts)]);
            let mut env = HashMap::new();
            env.insert("EMAIL".to_string(), "user@localhost".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
        }

        #[test]
        fn empty_local_part_fails() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Email);
            let schema = create_schema(vec![create_entry("EMAIL", opts)]);
            let mut env = HashMap::new();
            env.insert("EMAIL".to_string(), "@example.com".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
        }

        #[test]
        fn multiple_at_signs_fails() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Email);
            let schema = create_schema(vec![create_entry("EMAIL", opts)]);
            let mut env = HashMap::new();
            env.insert("EMAIL".to_string(), "user@@example.com".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
        }
    }

    // -------------------- Bool Kind Tests --------------------

    mod bool_kind_tests {
        use super::*;

        #[test]
        fn true_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Bool);
            let schema = create_schema(vec![create_entry("FLAG", opts)]);
            let mut env = HashMap::new();
            env.insert("FLAG".to_string(), "true".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn false_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Bool);
            let schema = create_schema(vec![create_entry("FLAG", opts)]);
            let mut env = HashMap::new();
            env.insert("FLAG".to_string(), "false".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn one_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Bool);
            let schema = create_schema(vec![create_entry("FLAG", opts)]);
            let mut env = HashMap::new();
            env.insert("FLAG".to_string(), "1".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn zero_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Bool);
            let schema = create_schema(vec![create_entry("FLAG", opts)]);
            let mut env = HashMap::new();
            env.insert("FLAG".to_string(), "0".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn yes_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Bool);
            let schema = create_schema(vec![create_entry("FLAG", opts)]);
            let mut env = HashMap::new();
            env.insert("FLAG".to_string(), "yes".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn no_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Bool);
            let schema = create_schema(vec![create_entry("FLAG", opts)]);
            let mut env = HashMap::new();
            env.insert("FLAG".to_string(), "no".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn case_insensitive() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Bool);
            let schema = create_schema(vec![create_entry("FLAG", opts)]);
            let mut env = HashMap::new();
            env.insert("FLAG".to_string(), "TRUE".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn invalid_bool_fails() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Bool);
            let schema = create_schema(vec![create_entry("FLAG", opts)]);
            let mut env = HashMap::new();
            env.insert("FLAG".to_string(), "maybe".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert!(errors[0].message.contains("not a valid boolean"));
        }
    }

    // -------------------- IP Address Kind Tests --------------------

    mod ip_address_kind_tests {
        use super::*;

        #[test]
        fn valid_ipv4_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::IpAddress);
            let schema = create_schema(vec![create_entry("IP", opts)]);
            let mut env = HashMap::new();
            env.insert("IP".to_string(), "192.168.1.1".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn localhost_ipv4_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::IpAddress);
            let schema = create_schema(vec![create_entry("IP", opts)]);
            let mut env = HashMap::new();
            env.insert("IP".to_string(), "127.0.0.1".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn valid_ipv6_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::IpAddress);
            let schema = create_schema(vec![create_entry("IP", opts)]);
            let mut env = HashMap::new();
            env.insert("IP".to_string(), "::1".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn full_ipv6_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::IpAddress);
            let schema = create_schema(vec![create_entry("IP", opts)]);
            let mut env = HashMap::new();
            env.insert(
                "IP".to_string(),
                "2001:0db8:85a3:0000:0000:8a2e:0370:7334".to_string(),
            );
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn invalid_ip_fails() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::IpAddress);
            let schema = create_schema(vec![create_entry("IP", opts)]);
            let mut env = HashMap::new();
            env.insert("IP".to_string(), "not.an.ip".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
            assert!(errors[0].message.contains("not a valid IP address"));
        }

        #[test]
        fn ip_out_of_range_fails() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::IpAddress);
            let schema = create_schema(vec![create_entry("IP", opts)]);
            let mut env = HashMap::new();
            env.insert("IP".to_string(), "256.256.256.256".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
        }

        #[test]
        fn hostname_fails() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::IpAddress);
            let schema = create_schema(vec![create_entry("IP", opts)]);
            let mut env = HashMap::new();
            env.insert("IP".to_string(), "localhost".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
        }
    }

    // -------------------- Path Kind Tests --------------------

    mod path_kind_tests {
        use super::*;

        #[test]
        fn absolute_path_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Path);
            let schema = create_schema(vec![create_entry("PATH", opts)]);
            let mut env = HashMap::new();
            env.insert("PATH".to_string(), "/usr/local/bin".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn relative_path_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Path);
            let schema = create_schema(vec![create_entry("PATH", opts)]);
            let mut env = HashMap::new();
            env.insert("PATH".to_string(), "./config.yaml".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn windows_path_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Path);
            let schema = create_schema(vec![create_entry("PATH", opts)]);
            let mut env = HashMap::new();
            env.insert("PATH".to_string(), "C:\\Users\\test".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn single_dot_passes() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Path);
            let schema = create_schema(vec![create_entry("PATH", opts)]);
            let mut env = HashMap::new();
            env.insert("PATH".to_string(), ".".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }
    }

    // -------------------- Combined Validation Tests --------------------

    mod combined_tests {
        use super::*;

        #[test]
        fn required_with_length_constraint() {
            let mut opts = default_options();
            opts.required = true;
            opts.required_length = Some(10);
            let schema = create_schema(vec![create_entry("API_KEY", opts)]);
            let mut env = HashMap::new();
            env.insert("API_KEY".to_string(), "1234567890".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn required_missing_skips_other_validations() {
            let mut opts = default_options();
            opts.required = true;
            opts.required_length = Some(10);
            opts.kind = Some(EntryKind::Integer {
                min: None,
                max: None,
            });
            let schema = create_schema(vec![create_entry("VAR", opts)]);
            let env = HashMap::new();
            let errors = validate_env(&schema, &env);
            // Should only have the "required" error, not length or kind errors
            assert_eq!(errors.len(), 1);
            assert!(errors[0].message.contains("missing"));
        }

        #[test]
        fn multiple_validation_errors_for_same_field() {
            let mut opts = default_options();
            opts.required_length = Some(5);
            opts.legal_values = Some(vec!["valid".to_string()]);
            let schema = create_schema(vec![create_entry("VAR", opts)]);
            let mut env = HashMap::new();
            env.insert("VAR".to_string(), "invalid_long".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 2);
        }

        #[test]
        fn multiple_entries_with_errors() {
            let mut opts1 = default_options();
            opts1.required = true;

            let mut opts2 = default_options();
            opts2.kind = Some(EntryKind::Integer {
                min: None,
                max: None,
            });

            let schema = create_schema(vec![
                create_entry("REQUIRED_VAR", opts1),
                create_entry("INT_VAR", opts2),
            ]);
            let mut env = HashMap::new();
            env.insert("INT_VAR".to_string(), "not_an_int".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 2);
        }

        #[test]
        fn integer_with_legal_values() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Integer {
                min: Some(1),
                max: Some(10),
            });
            opts.legal_values = Some(vec!["1".to_string(), "5".to_string(), "10".to_string()]);
            let schema = create_schema(vec![create_entry("LEVEL", opts)]);
            let mut env = HashMap::new();
            env.insert("LEVEL".to_string(), "5".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn regex_with_kind() {
            let mut opts = default_options();
            opts.regex_match = Some(r"^\d{3}$".to_string());
            opts.kind = Some(EntryKind::Integer {
                min: Some(100),
                max: Some(999),
            });
            let schema = create_schema(vec![create_entry("CODE", opts)]);
            let mut env = HashMap::new();
            env.insert("CODE".to_string(), "500".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }
    }

    // -------------------- Edge Cases --------------------

    mod edge_cases {
        use super::*;

        #[test]
        fn empty_schema() {
            let schema = create_schema(vec![]);
            let mut env = HashMap::new();
            env.insert("EXTRA".to_string(), "value".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn whitespace_only_value_treated_as_non_empty() {
            let mut opts = default_options();
            opts.required = true;
            let schema = create_schema(vec![create_entry("VAR", opts)]);
            let mut env = HashMap::new();
            env.insert("VAR".to_string(), "   ".to_string());
            let errors = validate_env(&schema, &env);
            // Whitespace-only is treated as non-empty
            assert!(errors.is_empty());
        }

        #[test]
        fn very_long_value() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::String);
            let schema = create_schema(vec![create_entry("LONG", opts)]);
            let mut env = HashMap::new();
            env.insert("LONG".to_string(), "x".repeat(100000));
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn special_characters_in_key() {
            let schema = create_schema(vec![create_entry("MY_VAR_123", default_options())]);
            let mut env = HashMap::new();
            env.insert("MY_VAR_123".to_string(), "value".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn newlines_in_value() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::String);
            let schema = create_schema(vec![create_entry("MULTI", opts)]);
            let mut env = HashMap::new();
            env.insert("MULTI".to_string(), "line1\nline2\nline3".to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn integer_at_i64_max() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Integer {
                min: None,
                max: None,
            });
            let schema = create_schema(vec![create_entry("BIG", opts)]);
            let mut env = HashMap::new();
            env.insert("BIG".to_string(), i64::MAX.to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn integer_at_i64_min() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Integer {
                min: None,
                max: None,
            });
            let schema = create_schema(vec![create_entry("SMALL", opts)]);
            let mut env = HashMap::new();
            env.insert("SMALL".to_string(), i64::MIN.to_string());
            let errors = validate_env(&schema, &env);
            assert!(errors.is_empty());
        }

        #[test]
        fn integer_overflow_fails() {
            let mut opts = default_options();
            opts.kind = Some(EntryKind::Integer {
                min: None,
                max: None,
            });
            let schema = create_schema(vec![create_entry("OVERFLOW", opts)]);
            let mut env = HashMap::new();
            env.insert("OVERFLOW".to_string(), "99999999999999999999".to_string());
            let errors = validate_env(&schema, &env);
            assert_eq!(errors.len(), 1);
        }
    }
}

mod kind_deserialization_tests {
    use crate::schema::{Entry, EntryKind};

    #[test]
    fn test_kind_deserialization() {
        let cases = vec![
            (
                "Integer",
                EntryKind::Integer {
                    min: None,
                    max: None,
                },
            ),
            (
                "integer",
                EntryKind::Integer {
                    min: None,
                    max: None,
                },
            ),
            ("ip_address", EntryKind::IpAddress),
            ("IP_ADDRESS", EntryKind::IpAddress),
            ("ipaddress", EntryKind::IpAddress),
        ];

        for (kind, ttype) in cases {
            let toml_str = format!(
                r#"
                key = "TEST_VAR"
                required = true
                kind = "{}"
            "#,
                kind
            );
            let entry: Entry = toml::from_str(&toml_str).expect("Failed to parse TOML");
            assert_eq!(entry.kind, Some(ttype));
        }
    }
}
