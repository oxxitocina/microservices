import mongoose from "mongoose"

const MovieData = new mongoose.Schema({
    kinopoiskId: {type: Number, required: true},
    imdbId: {type: String},
    nameRu: {type: String},
    nameEn: {type: String},
    nameOriginal: {type: String},
    countries: [{country: String}],
    genres: [{genre: String}],
    ratingKinopoisk: {type: Number},
    ratingImdb: {type: Number},
    year: {type: Number},
    type: {type: String},
    posterUrl: {type: String},
    posterUrlPreview: {type: String},
})

export default mongoose.model('MovieData', MovieData)