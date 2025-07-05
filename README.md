<<<<<<< HEAD
#
<h2>Fork Information</h2>

<p>This fork adds <strong>video thumbnail generation</strong> to FileBrowser.</p>

<p>Inspired by <a href="https://github.com/cody82/filebrowser/tree/videopreview">this fork</a> by cody82, but further customized.</p>

<h3>Customizations include:</h3>
<ul>
  <li>Thumbnails are stored in a <code>.thumbnails</code> subfolder, alongside the existing cache mechanism.</li>
  <li><strong>Hardware acceleration</strong> is dynamically selected based on the OS for maximum compatibility:
    <ul>
      <li><strong>Linux</strong>: VA-API</li>
      <li><strong>Windows</strong>: DXVA2</li>
      <li><strong>Mac</strong>: VideoToolbox</li>
    </ul>
  </li>
  <li>Thumbnail generation <strong>pauses when leaving a folder</strong>, allowing new folders to start processing without waiting for previous ones to finish.</li>
</ul>

<h3>Requirements</h3>
<ul>
  <li><strong>FFmpeg</strong> is required.</li>
  <li>The <strong>ctx utility</strong> is used to handle stopping thumbnail generation.</li>
</ul>

<h3>Notes</h3>
<ul>
  <li><strong>22-Feb-25</strong> – This fork will now be <strong>automatically kept up to date</strong> via GitHub Actions.</li>
  <li>The following files were modified:
    <ul>
      <li><code>/http/preview.go</code></li>
      <li><code>/frontend/src/components/files/ListingItem.vue</code></li>
    </ul>
  </li>
</ul>

<h3>About This Project</h3>
<p>ChatGPT was heavily involved in generating the code.</p>
<p>Although I have programming experience and even a <strong>software engineering degree</strong>, my career path led me to <strong>IT and System Administration</strong>.</p>
<p>This is my <strong>first real experience with Golang</strong>, so go easy on me! 😅</p>


=======
>>>>>>> upstream/master
<p align="center">
  <img src="https://raw.githubusercontent.com/filebrowser/logo/master/banner.png" width="550"/>
</p>

[![Build](https://github.com/filebrowser/filebrowser/actions/workflows/main.yaml/badge.svg)](https://github.com/filebrowser/filebrowser/actions/workflows/main.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/filebrowser/filebrowser)](https://goreportcard.com/report/github.com/filebrowser/filebrowser)
[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/filebrowser/filebrowser)
[![Version](https://img.shields.io/github/release/filebrowser/filebrowser.svg)](https://github.com/filebrowser/filebrowser/releases/latest)
[![Chat IRC](https://img.shields.io/badge/freenode-%23filebrowser-blue.svg)](http://webchat.freenode.net/?channels=%23filebrowser)

File Browser provides a file managing interface within a specified directory and it can be used to upload, delete, preview and edit your files. It is a **create-your-own-cloud**-kind of software where you can just install it on your server, direct it to a path and access your files through a nice web interface.

## Documentation

Documentation on how to install, configure, and contribute to this project is hosted at [filebrowser.org](https://filebrowser.org).

## Project Status

> [!WARNING]
>
> This project is currently on **maintenance-only** mode, and is looking for new maintainers. For more information, please read the [discussion #4906](https://github.com/filebrowser/filebrowser/discussions/4906). Therefore, please note the following:
>
> - It can take a while until someone gets back to you. Please be patient.
> - [Issues][issues] are only being used to track bugs. Any unrelated issues will be converted into a [discussion][discussions].
> - No new features will be implemented until further notice. The priority is on triaging issues and merge bug fixes.
> 
> If you're interested in maintaining this project, please reach out via the discussion above.

[issues]: https://github.com/filebrowser/filebrowser/issues
[discussions]: https://github.com/filebrowser/filebrowser/discussions

## Contributing

Contributions are always welcome. To start contributing to this project, read our [guidelines](CONTRIBUTING.md) first.

## License

[Apache License 2.0](LICENSE) © File Browser Contributors
