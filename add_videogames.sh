echo "{ \"title\": \"Super Mario Odyssey\", \"releaseDate\": \"2017-10-27\", \"publisher\": \"Nintendo\" }" | curl -d @- -X POST http://localhost:8080/videogame/create

echo "{ \"title\": \"Super Mario Galaxy\", \"releaseDate\": \"2007-11-01\", \"publisher\": \"Nintendo\" }" | curl -d @- -X POST http://localhost:8080/videogame/create

echo "{ \"title\": \"Halo 3\", \"releaseDate\": \"2007-09-25\", \"publisher\": \"Xbox Game Studios\" }" | curl -d @- -X POST http://localhost:8080/videogame/create

echo "{ \"title\": \"Crash Bandicoot 4: It's About Time\", \"releaseDate\": \"2020-10-02\", \"publisher\": \"Activision\" }" | curl -d @- -X POST http://localhost:8080/videogame/create

echo "{ \"title\": \"Edna & Harvey: The Breakout\", \"releaseDate\": \"2008-06-05\", \"publisher\": \"Daedalic Entertainment\" }" | curl -d @- -X POST http://localhost:8080/videogame/create

echo "{ \"title\": \"God of War\", \"releaseDate\": \"2018-04-20\", \"publisher\": \"Sony Interactive Entertainment\" }" | curl -d @- -X POST http://localhost:8080/videogame/create

echo "{ \"title\": \"The Elder Scrolls V: Skyrim\", \"releaseDate\": \"2011-11-11\", \"publisher\": \"Bethesda Softworks\" }" | curl -d @- -X POST http://localhost:8080/videogame/create