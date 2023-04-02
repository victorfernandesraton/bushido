use reqwest::blocking;
use reqwest::Error as RequestError;
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

#[derive(Debug, Serialize, Deserialize)]
struct MangaData {
    id_serie: u64,
    label: String,
    link: String,
}

#[derive(Deserialize, Debug)]
#[serde(untagged)]
enum MangaSearchSerie {
    Simple(bool),
    Explicit(Vec<MangaData>)
}

#[derive(Deserialize, Debug)]
pub struct MangaSearchResponse {
    series: MangaSearchSerie
}


fn get_manga_images_reqyest(id: u64) -> Result<MangaPagesResponse, RequestError> {
    let url = format!("https://mangalivre.net/leitor/pages/{}", id);
    let resposta = blocking::get(url)?.json::<MangaPagesResponse>()?;
    Ok(resposta)
}

pub fn search(query: &str) -> Result<MangaSearchResponse, RequestError> {
    let client = blocking::Client::new();
    let body = [("search", query)];

    let response = client
        .post("https://mangalivre.net/lib/search/series.json")
        .form(&body)
        .header("X-Requested-With", "XMLHttpRequest")
        .send()
        .unwrap()
        .json()?;

    Ok(response)
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
