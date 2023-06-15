import path from 'path'
import * as grpc from '@grpc/grpc-js'
import * as protoLoader from '@grpc/proto-loader'
import {ProtoGrpcType} from './proto/random'
import readline from 'readline'
import amqp from 'amqplib/callback_api.js'
import MovieVideo from './movieVideo'
import mongoose from 'mongoose'

const PORT = 8082
const PROTO_FILE = './proto/random.proto'
const DB_URL = 'mongodb+srv://user:user@cluster0.4b9gqtt.mongodb.net/'

const packageDef = protoLoader.loadSync(path.resolve(__dirname, PROTO_FILE))
const grpcObj = (grpc.loadPackageDefinition(packageDef) as unknown) as ProtoGrpcType


const client = new grpcObj.randomPackage.Random(
  `0.0.0.0:${PORT}`, grpc.credentials.createInsecure()
)

const deadline = new Date()
deadline.setSeconds(deadline.getSeconds() + 5)
client.waitForReady(deadline, (err) => {
  if (err) {
    console.error(err)
    return
  }
  onClientReady()
})


function onClientReady() {
  client.PingPong({message: "Ping"}, (err, result) => {
    if (err) {
      console.error(err)
      return
    }
    console.log(result)
  })

  amqp.connect('amqps://bzecpdsx:ohp1UJe_qnSflyW65NGAPiqTKhoDQoxY@hawk.rmq.cloudamqp.com/bzecpdsx', (error0: any, connection: any) => {
    if(error0) {
        throw error0
    }

    connection.createChannel((error1: any, channel: any) => {
        if(error1) {
            throw error1
        }

        mongoose.connect(DB_URL);
        
        channel.assertQueue('getVideos', {durable: false})
        channel.consume('getVideos', (msg: any) => {
            try {
                const eventObject = JSON.parse(msg.content.toString())
                console.log(eventObject)
                
                const metadata = new grpc.Metadata()
                metadata.set("username", username)
                const call = client.Chat(metadata)
                call.write({
                  message: eventObject
                })

                call.on("data", async (chunk) => {
                  console.log(chunk.message)
                  const data = JSON.parse(chunk.message)
                  
                  await MovieVideo.create(data)
                  console.log('video links created')
                })

            } catch (error) {
                console.log(error)
            }
        })





        // channel.assertQueue('genreId', {durable: false})
        // channel.consume('genreId', (msg: any) => {
        //     try {
        //         // const obj = Post.decode(msg.content)
        //         // console.log(obj)
        //         const eventObject = JSON.parse(msg.content.toString())
        //         console.log(eventObject)
                
        //         const metadata = new grpc.Metadata()
        //         metadata.set("username", username)
        //         const call = client.Chat(metadata)
        //         call.write({
        //           message: eventObject
        //         })

        //         call.on("data", (chunk) => {
        //           // console.log(`${chunk.username} ==> ${chunk.message}`)
        //           console.log(chunk.message)
        //           channel.sendToQueue('movies', Buffer.from(chunk.message))
        //           // channel.sendToQueue('movieId', Buffer.from(JSON.stringify(eventObject)))
        //         })

        //     } catch (error) {
        //         console.log(error)
        //     }
        // })

        

        




    })
})




  const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
  })

  const username = process.argv[2]
  if (!username) console.error("No username, can't join chat"), process.exit()


  const metadata = new grpc.Metadata()
  metadata.set("username", username)
  const call = client.Chat(metadata)

  call.on("data", (chunk) => {
    // console.log(`${chunk.username} ==> ${chunk.message}`)
    console.log(chunk.message)
  })

  // rl.on("line", (line) => {
  //   if(line === "quit") {
  //     call.end()
  //   }else {
  //     call.write({
  //       message: `${line} testtsettest`
  //     })
  //   }

  // })
}

