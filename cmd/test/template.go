package main

const tpl = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<script src="https://cdn.tailwindcss.com"></script>
<title>Image Grid</title>
</head>
<body class="bg-gray-900 text-white p-4">

<h1 class="text-3xl font-bold mb-6">Gallery</h1>

<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
{{range .}}
  <div class="bg-gray-800 rounded-lg overflow-hidden shadow-lg">
    <a href="{{.Link}}">
      <img src="{{.Image}}" alt="{{.TitleH2}}" class="w-full h-64 object-cover hover:scale-105 transition-transform duration-300">
    </a>
    <div class="p-4">
      <h5 class="text-sm font-semibold text-gray-300">{{.TitleH5}}</h5>
      <h2 class="text-lg font-bold">{{.TitleH2}}</h2>
    </div>
  </div>
{{end}}
</div>

</body>
</html>
`
const tpl2 = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<script src="https://cdn.tailwindcss.com"></script>
<title>{{.Title}}</title>
</head>
<body class="bg-gray-900 text-white p-4">

<h1 class="text-3xl font-bold mb-6">{{.Title}}</h1>

<div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
{{range .Images}}
  <div class="bg-gray-800 rounded-lg overflow-hidden shadow-lg">
    <a href="{{.Link}}" download>
      <img src="{{.Link}}" class="w-full h-64 object-cover hover:scale-105 transition-transform duration-300">
    </a>
  </div>
{{end}}
</div>

</body>
</html>
`
