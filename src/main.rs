mod api;

fn main() {
    let result = api::search("solo");
    print!("{:?}", result)
}
