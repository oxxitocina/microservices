import express from 'express'
import mongoose from 'mongoose'
import cors from 'cors'
import amqp from 'amqplib/callback_api.js'
import MovieData from './movieData.js';
import MovieVideo from './movieVideo.js';
import protobuf from 'protobufjs'

const PORT = 5000;
const DB_URL = 'mongodb+srv://user:user@cluster0.4b9gqtt.mongodb.net/'

const app = express();
app.use(express.json())
app.use(cors({
    origin: ['http://localhost:1234', 'http://localhost:5000','http://localhost:3000',]
}))

async function startApp()   {
    try {
        const root = await protobuf.load('movieId.proto')
        const MovieIdProto = root.lookupType('movieidpackage.MovieId')

        amqp.connect('amqps://bzecpdsx:ohp1UJe_qnSflyW65NGAPiqTKhoDQoxY@hawk.rmq.cloudamqp.com/bzecpdsx', (error0, connection) => {
            if(error0) {
                throw error0
            }
            connection.createChannel((error1, channel) => {
                if(error1) {
                    throw error1
                }
                mongoose.connect(DB_URL);
 
                app.post('/api/movies', async(req, res) => {
                    try {   
                        const { kinopoiskId } = req.body
                        const encodedMovie = MovieIdProto.encode({id: kinopoiskId}).finish()
                        channel.sendToQueue('postMovie', Buffer.from(JSON.stringify(req.body)))
                        channel.sendToQueue('getVideos', Buffer.from(encodedMovie))
                        return res.status(200).json({status: 'delivered'});
                    } catch (error) {
                        console.log(error)
                        res.status(500).json(error)
                    }
                })

                app.get('/api/movies/:id', async(req, res) => {
                    try {
                        const movie = await MovieData.findOne({'kinopoiskId': req.params.id})
                        const movieVid = await MovieVideo.findOne({'id': req.params.id})
                        const obj = [movie, movieVid]
                        return res.json(obj)

                    } catch (error) {
                        res.status(500).json(error)
                    }
                })      

                app.listen(PORT, () => console.log('server started on port ' + 5000));
                process.on('beforeExit', () => {
                    console.log('closing')
                    connection.close();
                })

            })
        })

        return 'successful connection'

    } catch (error) {
        console.log(error)
    }
}

startApp();

test('test', async () => {
    const data = await startApp();
    expect(data).toBe('successful connection');
  });