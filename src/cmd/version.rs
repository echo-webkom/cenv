use reqwest::blocking::Client;
use serde::Deserialize;

pub const VERSION: Option<&'static str> = option_env!("CENV_VERSION");

#[derive(Debug, Deserialize)]
struct VersionInfo {
    tag_name: String,
}

pub fn check_latest_version_and_warn() {
    let _ = try_check_latest_version();
}

// Check the latest version of cenv on GitHub and print a warning if the current version is
// outdated.
fn try_check_latest_version() -> Option<()> {
    let version = VERSION?;

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

    Some(())
}
