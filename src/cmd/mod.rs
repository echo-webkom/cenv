use clap::{Parser, Subcommand};
use std::collections::HashMap;
use std::fs;
use std::path::PathBuf;

use crate::schema::{Schema, generate_env, validate_env};

#[derive(Parser)]
#[command(name = "cenv")]
#[command(about = "Environment file manager using schema definitions", long_about = None)]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// Generate .env or add missing values from the schema
    Fix {
        /// Path to the schema file
        #[arg(short, long, default_value = "cenv.schema.toml")]
        schema: PathBuf,

        /// Path to the .env file
        #[arg(short, long, default_value = ".env")]
        env: PathBuf,
    },
    /// Validate .env against the schema
    Check {
        /// Path to the schema file
        #[arg(short, long, default_value = "cenv.schema.toml")]
        schema: PathBuf,

        /// Path to the .env file
        #[arg(short, long, default_value = ".env")]
        env: PathBuf,
    },
}

pub fn run() {
    let cli = Cli::parse();

    match cli.command {
        Commands::Fix { schema, env } => {
            if let Err(e) = fix_command(&schema, &env) {
                eprintln!("Error: {}", e);
                std::process::exit(1);
            }
        }
        Commands::Check { schema, env } => {
            if let Err(e) = check_command(&schema, &env) {
                eprintln!("Error: {}", e);
                std::process::exit(1);
            }
        }
    }
}

fn fix_command(
    schema_path: &PathBuf,
    env_path: &PathBuf,
) -> Result<(), Box<dyn std::error::Error>> {
    let schema_contents = fs::read_to_string(schema_path)?;
    let schema: Schema = toml::from_str(&schema_contents)?;

    let existing_env: HashMap<String, String> = if env_path.exists() {
        dotenvy::from_path_iter(env_path)?
            .filter_map(|item| item.ok())
            .collect()
    } else {
        HashMap::new()
    };

    // Generate and write .env content
    let env_content = generate_env(&schema, &existing_env);
    fs::write(env_path, env_content)?;

    println!("Successfully updated {}", env_path.display());
    Ok(())
}

fn check_command(
    schema_path: &PathBuf,
    env_path: &PathBuf,
) -> Result<(), Box<dyn std::error::Error>> {
    let schema_contents = fs::read_to_string(schema_path)?;
    let schema: Schema = toml::from_str(&schema_contents)?;

    if !env_path.exists() {
        return Err(".env file does not exist".into());
    }
    let env_vars: HashMap<String, String> = dotenvy::from_path_iter(env_path)?
        .filter_map(|item| item.ok())
        .collect();

    let errors = validate_env(&schema, &env_vars);

    if errors.is_empty() {
        Ok(())
    } else {
        for error in &errors {
            eprintln!("{}", error);
        }
        std::process::exit(1);
    }
}
