import mongoose from "mongoose"

const MovieVideo = new mongoose.Schema({
    id: { type: Number },
    total: {type: String},
    items: [{url: String}, {name: String}, {site: String}]
})

export default mongoose.model('MovieVideo', MovieVideo)