mod api;

fn main() {
    let images =
        api::get_images("https://mangalivre.net/ler/solo-leveling/online/453669/191#/!page0");
    print!("{:?}", images)
}
