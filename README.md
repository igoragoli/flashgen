# flashgen

A very simple CLI tool that turns plaintext documents into Anki flashcards.

## Installation

Make sure you have [Go 1.21 or later installed](https://go.dev/doc/install) in your machine.

```bash
go version
```

Install the package with

```bash
go install github.com/igoragoli/flashgen
```

Notes that have code or math are formatted to be compatible with the [Markdown and KaTeX add-on](https://ankiweb.net/shared/info/1087328706). Be sure to get it if you're going to need these kinds of notes.

## Usage

Initially, export your OpenAI key environment variable. Set it up [here](https://platform.openai.com/docs/quickstart/step-2-setup-your-api-key) if you don't have one.

```bash
export OPEN_AI_KEY=YOUR_OPEN_AI_KEY
```

To generate flashcards from a plaintext file, run:

```bash
flashgen <input-file> [<output-file>]
```

Flashcards will be saved in a .csv format interpretable by [Anki's import feature](https://docs.ankiweb.net/importing/intro.html). You can open the output file in your spreadsheet editor of choice to edit and rewrite cards before importing them to Anki.

If you're sure that you'd like to import your cards to Anki right away, run:

```bash
flashgen -anki -deck "<deck-name>" <input-file> [<output-file>]
```

If `-deck` is not provided, a new deck will be created with an appropriate name. `<deck-name>` will be created if it does not exist.

By default, the generated flashcards won't have tags set to them, because tags are applied in widely different ways by Anki users. If you'd like to set tags for all generated flashcards, run:

```bash
flashgen -anki -deck "<deck-name>" -tags "<tag1> <tag2>" <input-file> [<output-file>]
```

Tags have to be separated by whitespaces, as in

```bash
flashgen -anki -deck "go" -tags "go::functions go::concurrency" go-notes.md go-flashcards.csv
```

## Roadmap

- [x] Base functionality
- [ ] Support for tags
- [ ] Automatically adding

## FAQ

### What about PDF files?

As of now, there is no plan to support PDF files. As a workaround, you can use one of the following tools to convert your PDF file to a plaintext one:

- [pdftotext](https://www.xpdfreader.com/pdftotext-man.html)
- [nougat](https://github.com/facebookresearch/nougat)
- [Mathpix](https://mathpix.com/) (paid)
- Online PDF OCR tools, such as [this one](https://tools.pdf24.org/en/ocr-pdf)

### But writing your own questions is super important when using Anki!

Yes, I know, and as someone who uses Anki daily for years, I fully agree! 

I've built this to solve a very particular need: I always add notes to Anki when I'm studying a technical topic. Usually I 1. Read a chapter/article; 2. Write flashcards for it; 3. Go to step 1. 

But I constantly procrastinated advancing into some books or articles because I was lagging behind on my flashcards, so I needed a way to create flashcards easily.