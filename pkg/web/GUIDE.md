# web package

`web` is a package that provides a framework for building web applications in Go.
It provides some tools to define **static** web applications.

## Concepts

- **App**
  - The root of the web application.
  - Each app is mounted to a specific path, and its pages has path with prefix of the app's path.
  - App can have some static files and stylesheets, which are available for all pages in the app.
- **Page**
  - Page is a HTML web page.
  - Each path has a path.
  - Page can have its own stylesheets.
  - Page contains 'root component', which will be inserted into the `<body>` of the HTML document.
  - There are some objective of the page handler.
    1. From the HTTP request, process all required tasks and prepare data for the HTML response.

- **Component**
  - In `web` package, a component is a builder
