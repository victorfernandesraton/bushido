# Bushido: A Samurai Path to Manage Your Manga

## Origin

This project was created as an alternative solution to the [Komikku](https://valos.gitlab.io/Komikku/) project, a manga reader app for Linux. While the Komikku project is awesome, I found that it had some issues with how things were done.

### Distribution Problem

[Komikku](https://valos.gitlab.io/Komikku/) is developed using Python wrappers as a builder project, an integrated IDE for creating GNOME projects. This is an issue if you want to make it a multiplatform project.

While Builder is a great solution for delivering flatpack-based apps for Linux desktop platforms, it's not possible to ship it for Windows or Mac.

To address this issue, Bushido was created as a multiplatform solution for managing manga collections.

### Detach Problem

If you're wondering why I didn't use Flutter or web-based tools like Electron or Tauri, it's not viable either. For me, if you make a great CLI tool, someone might be interested in creating a Flutter client for you.

And for web-based things, using something based on a client-server approach like Electron seems to have issues with CORS and using third-party content like images inside a CDN.

### Why Go, Not Rust or JS-Thing?

For things like this, Rust might have better performance, but handling things like web scraping in Rust seems tough. I'm not interested in making a blazing-fast thing. As for JavaScript, I just don't want to use the chaotic JS things for a non-ideal solution.
