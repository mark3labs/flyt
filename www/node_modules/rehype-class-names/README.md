# Rehype Class Names

Adding classes to elements with rehype. Includes typescript support

## Installation

```
npm install rehype-class-names
```

## API

### `rehype().use(addClasses, additions])`

Add to `rehype` or `unified` pipeline with `.use`, where `additions` is an object
with keys that are the css selectors and the values are the classes you want to add.  

For more information on what keys you can pass, see the [`hast-util-select`](https://github.com/syntax-tree/hast-util-select#support) documentation.
For more information on what values you can pass, see the [`hast-util-classnames`](https://github.com/syntax-tree/hast-util-classnames#use) documentation.

Example:

```js
.use(addClasses, {
    'h1,h2,h3': 'title',
    h1: 'is-1',
    h2: 'is-2',
    h3: { 'hello': false, 32: true },
    h4: 10,
    h5: [],
    p: 'one two',
    a: "test",
    ol: "list-decimal",
    li: ['test', 'mega'],
})
```

This will output

```html
<h1 class="title is-1">header</h1>
<h2 class="title is-2">sub 2</h2>
<h3 class="title 32">sub 3</h3>
<h4 class="10">sub 4</h4>
<h5 class="">Sub 5</h5>
<p class="one two">
  Profile pictures are important.
  <a href="riderjensen.com" class="test">Rider Jensen</a>
</p>
<ol class="list-decimal">
  <li class="test mega">Create a Github repository</li>
  <li class="test mega">Pushed up code</li>
  <li class="test mega">Connect account to Github</li>
</ol>
```