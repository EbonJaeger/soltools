use clap::{Parser, Subcommand};
use thiserror::Error;

mod init;

#[derive(Parser)]
#[command(author, version, long_about = None, about = "Tool to make Solus packaging even easier",)]
struct Cli {
    #[command(subcommand)]
    command: Option<Commands>,
}

#[derive(Subcommand)]
enum Commands {
    /// Initialize a new package repository
    Init {
        /// Name of the package
        name: String,
        /// URL to the source tarball
        url: String,
        /// Create a maintainers file for this package
        #[arg(short, long)]
        maintain: bool,
    },
}

pub fn process() -> Result<(), Error> {
    let cli = Cli::parse();

    match &cli.command {
        Some(Commands::Init {
            name,
            url,
            maintain,
        }) => {
            init::handle(name.to_string(), url.to_string(), *maintain)?;
        }
        None => {}
    }

    Ok(())
}

#[derive(Debug, Error)]
pub enum Error {
    #[error("Package init error")]
    Init(#[from] init::Error),
}
