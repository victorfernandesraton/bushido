use reqwest::blocking;
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
struct MangaPage {
    legacy: String,
    avif: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct MangaPagesResponse {
    images: Vec<MangaPage>,
}

fn get_manga_images_reqyest(id: u64) -> Result<MangaPagesResponse, reqwest::Error> {
    let url = format!("https://mangalivre.net/leitor/pages/{}", id);
    let resposta = blocking::get(url)?.json::<MangaPagesResponse>()?;
    Ok(resposta)
}

pub fn get_images(url: &str) -> Result<Vec<String>, Box<dyn std::error::Error>> {
    let elementos: Vec<&str> = url.split('/').collect();
    let id = elementos
        .get(6)
        .ok_or("Unable to parse ID")?
        .parse::<u64>()?;
    let response = get_manga_images_reqyest(id)?;
    let items = response
        .images
        .iter()
        .map(|manga_page| manga_page.legacy.clone())
        .collect();
    Ok(items)
}
