package main

const frontpage_etag = "ad2cc2bf9b0d115e067c4d8b799f0852833f1747"

const frontpage = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" http-equiv="X-UA-Compatible" content="IE=edge;">
    <title id='title'>File Drop</title>
  </head>
  <body>
    <form enctype="multipart/form-data" action="/" method="post">
      Sender Email:<br>
      <input type="text" name="sender"><br>
      Recipients:<br>
      <input type="text" name="recipients"></br>
      Filename: <br>
      <input type="file" name="filename"></br>
      Duration: <br>
      <select name="factor">
	<option value="hour">hours</option>
	<option value="day">days</option>
	<option value="week">weeks</option>
	<option value="month">months</option>
      </select>
      <input type="number" name="duration" min="1" max="24"><br>
      <input type="submit" value="Submit">
    </form>
  </body>
</html>
`
