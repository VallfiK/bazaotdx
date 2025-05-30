const mongoose = require('mongoose');

const cottageSchema = new mongoose.Schema({
    name: {
        type: String,
        required: true
    },
    description: {
        type: String,
        required: true
    },
    capacity: {
        type: Number,
        required: true
    },
    price: {
        type: Number,
        required: true
    },
    images: [{
        type: String
    }],
    status: {
        type: String,
        enum: ['free', 'booked', 'bought'],
        default: 'free'
    },
    amenities: [{
        type: String
    }],
    createdAt: {
        type: Date,
        default: Date.now
    }
});

module.exports = mongoose.model('Cottage', cottageSchema);
