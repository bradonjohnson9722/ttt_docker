db = db.getSiblingDB('tic-tac-toe');  // Switch to or create the database
db.createCollection('game');  // Create a collection
db.mycollection.insert({ name: "example" });  // Insert example data
