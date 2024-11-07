db.message.createIndex( { userId: 1, createdAt: -1 } )
db.user.createIndex( { username: 1 }, { unique: true } )

# create cloudinary
db("luce").createCollection("images")
