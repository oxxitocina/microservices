import mongoose from 'mongoose'

const DB_URL = 'mongodb+srv://user:user@cluster0.4b9gqtt.mongodb.net/'

async function connectToMongo()   {
    try {
        mongoose.connect(DB_URL);
        return 'connected to mongo'
    } catch (error) {
        return error
    }
}

test('test', async () => {
    const data = await connectToMongo();
    expect(data).toBe('connected to mongo');
});