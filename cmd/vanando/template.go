package main

import "html/template"

var tpl = template.Must(template.New("gallery").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Image Gallery</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            padding: 20px;
            margin: 0;
            background: #f5f5f5;
        }
        h1, h2 {
            text-align: center;
        }
        form {
            display: flex;
            justify-content: center;
            margin-bottom: 20px;
            flex-wrap: wrap;
        }
        input[type="text"] {
            width: 70%;
            max-width: 500px;
            padding: 10px;
            margin: 5px;
            border: 1px solid #ccc;
            border-radius: 5px;
            box-sizing: border-box;
        }
        input[type="submit"] {
            padding: 10px 20px;
            margin: 5px;
            border: none;
            border-radius: 5px;
            background-color: #007bff;
            color: white;
            cursor: pointer;
        }
        input[type="submit"]:hover {
            background-color: #0056b3;
        }
        .gallery {
            display: flex;
            flex-wrap: wrap;
            justify-content: center;
            gap: 15px;
        }
        .gallery img {
            max-width: 100%;
            height: auto;
            border-radius: 10px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.2);
        }
        .gallery-item {
            flex: 1 1 200px;
            max-width: 300px;
        }
    </style>
</head>
<body>
    <h1>Paste Image URLs</h1>
    <form method="post">
        <input type="text" name="url" placeholder="Enter image URL" required>
        <input type="submit" value="Add">
    </form>

    <h2>Gallery</h2>
    {{if .}}
    <div class="gallery">
        {{range .}}
        <div class="gallery-item">
            <img src="{{.}}" alt="Image">
        </div>
        {{end}}
    </div>
    {{else}}
        <p style="text-align:center;">No images yet.</p>
    {{end}}
</body>
</html>
`))
