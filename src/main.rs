mod api;
use reqwest::Error as RequestError;


fn main() -> Result<(), RequestError> {
    let result: api::MangaSearchResponse = api::search("solo")?;
    let series = result.series;
    let mut series_list: Vec<api::MangaData> = Vec::new();
    match series {
        api::MangaSearchSerie::Explicit(series) => series_list = series,
        _ => {}
    }
    for serie in series_list {
        // print!("{:?}", serie);
        let chapters = api::get_chapters(serie.id_serie)?;
        let chapter_data = chapters.chapters;
        let mut chapter_list: Vec<api::MangaChapterData> = Vec::new();
        match chapter_data {
            api::MangaChapter::Explicit(chapter) => chapter_list = chapter,
            _ => {}
        }
        for ch in chapter_list  {
            let path = serie.link.split("/").collect::<Vec<&str>>()[2];
            let url = format!("https://mangalivre.net/ler/{}/online/{}/{}",path, serie.id_serie, ch.number);
            let images = api::get_images(&url);
            print!("{:?}",images);
        }

    }

    // print!("{:?}", images);

    Ok(())
}
