use std::error::Error;

mod cli;

fn main() -> Result<(), Box<dyn Error>> {
    cli::process()?;

    Ok(())
}
