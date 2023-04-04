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
pub struct MangaData {
    pub id_serie: u64,
    label: String,
    pub link: String,
}

#[derive(Deserialize, Debug)]
#[serde(untagged)]
pub enum MangaSearchSerie {
    Simple(bool),
    Explicit(Vec<MangaData>)
}


#[derive(Deserialize, Debug)]

pub struct MangaSearchResponse {
    pub series: MangaSearchSerie
}

#[derive(Deserialize, Debug)]

pub struct MangaChapterData {
    id_serie: u64,
    pub id_chapter: u64,
    name: String,
    chapter_name:String,
    pub number: String,
    // TODO: parse datetime
    date_created: String
}

#[derive(Deserialize, Debug)]
#[serde(untagged)]
pub enum MangaChapter {
    Simple(bool),
    Explicit(Vec<MangaChapterData>)
}

#[derive(Deserialize, Debug)]
pub struct MangaChapterResponse {
    pub chapters: MangaChapter
}




fn get_manga_images_reqyest(id: u64) -> Result<MangaPagesResponse, RequestError> {
    let url = format!("https://mangalivre.net/leitor/pages/{}", id);
    let resposta = blocking::get(url)?.json::<MangaPagesResponse>()?;
    Ok(resposta)
}

pub fn search(query: &str) -> Result<MangaSearchResponse, RequestError> {
    let client = blocking::Client::new();
    let body = [("search", query)];

    let response: MangaSearchResponse = client
        .post("https://mangalivre.net/lib/search/series.json")
        .form(&body)
        .header("X-Requested-With", "XMLHttpRequest")
        .send()
        .unwrap()
        .json()?;


    Ok(response)
}

fn get_manga_id_by_url(url: &str) -> Result<u64,  Box<dyn std::error::Error>> {
    let elementos: Vec<&str> = url.split('/').collect();
    let id = elementos
        .get(6)
        .ok_or("Unable to parse ID")?
        .parse::<u64>()?;
    Ok(id)
}

pub fn get_images(url: &str) -> Result<Vec<String>, Box<dyn std::error::Error>> {
    let id = get_manga_id_by_url(url)?;
    let response = get_manga_images_reqyest(id)?;
    let items = response
        .images
        .iter()
        .map(|manga_page| manga_page.legacy.clone())
        .collect();
    Ok(items)
}

pub fn get_chapters(id: u64) -> Result<MangaChapterResponse, RequestError> {
    let url = format!("https://mangalivre.net/series/chapters_list.json?page=1&id_serie={}", id);

    let client = blocking::Client::new();

    let response = client
        .get(url)
        .header("X-Requested-With", "XMLHttpRequest")
        .send()
        .unwrap().json::<MangaChapterResponse>()?;
    
    Ok(response)
}