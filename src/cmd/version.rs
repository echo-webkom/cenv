use reqwest::blocking::Client;
use serde::{Deserialize, Serialize};
use std::fs;
use std::path::PathBuf;
use std::time::{SystemTime, UNIX_EPOCH};

pub const VERSION: Option<&'static str> = option_env!("CENV_VERSION");

const CACHE_TTL_SECS: u64 = 86400; // 24 hours

#[derive(Debug, Deserialize)]
struct VersionInfo {
    tag_name: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct VersionCache {
    time: u64,
    tag_name: String,
}

pub fn check_latest_version_and_warn() {
    let _ = try_check_latest_version();
}

fn cache_path() -> PathBuf {
    std::env::temp_dir().join("cenv_version_check")
}

fn now_secs() -> u64 {
    SystemTime::now()
        .duration_since(UNIX_EPOCH)
        .map(|d| d.as_secs())
        .unwrap_or(0)
}

fn read_cache() -> Option<VersionCache> {
    let contents = fs::read_to_string(cache_path()).ok()?;
    serde_json::from_str(&contents).ok()
}

fn write_cache(info: &VersionInfo) {
    let cache = VersionCache {
        time: now_secs(),
        tag_name: info.tag_name.clone(),
    };
    if let Ok(json) = serde_json::to_string(&cache) {
        let _ = fs::write(cache_path(), json);
    }
}

// Check the latest version of cenv on GitHub and print a warning if the current version is
// outdated.
fn try_check_latest_version() -> Option<()> {
    let version = VERSION?;

    if let Some(cache) = read_cache()
        && now_secs().saturating_sub(cache.time) < CACHE_TTL_SECS
    {
        if cache.tag_name != version {
            eprintln!(
                "Warning: A new version of cenv is available: {} (current: {})",
                cache.tag_name, version
            );
            eprintln!();
        }
        return Some(());
    }

    let client = Client::new();
    let resp = client
        .get("https://api.github.com/repos/echo-webkom/cenv/releases/latest")
        .header("User-Agent", "cenv")
        .send()
        .ok()?;

    let info = resp.json::<VersionInfo>().ok()?;

    if info.tag_name != version {
        eprintln!(
            "Warning: A new version of cenv is available: {} (current: {})",
            info.tag_name, version
        );
        eprintln!();
    }

    write_cache(&info);

    Some(())
}
