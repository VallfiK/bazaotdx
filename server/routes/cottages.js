const express = require('express');
const router = express.Router();
const auth = require('../middleware/auth');
const Cottage = require('../models/Cottage');

// @route   GET /api/cottages
// @desc    Get all cottages
router.get('/', async (req, res) => {
    try {
        const cottages = await Cottage.find();
        res.json(cottages);
    } catch (err) {
        console.error(err.message);
        res.status(500).send('Server error');
    }
});

// @route   POST /api/cottages
// @desc    Add new cottage
// @access  Private
router.post('/', [auth], async (req, res) => {
    const { name, description, capacity, price, images, amenities } = req.body;

    try {
        const newCottage = new Cottage({
            name,
            description,
            capacity,
            price,
            images,
            amenities
        });

        const cottage = await newCottage.save();
        res.json(cottage);
    } catch (err) {
        console.error(err.message);
        res.status(500).send('Server error');
    }
});

// @route   PUT /api/cottages/:id
// @desc    Update cottage
// @access  Private
router.put('/:id', [auth], async (req, res) => {
    try {
        const cottage = await Cottage.findById(req.params.id);
        if (!cottage) {
            return res.status(404).json({ msg: 'Cottage not found' });
        }

        const updatedCottage = await Cottage.findByIdAndUpdate(
            req.params.id,
            { $set: req.body },
            { new: true }
        );

        res.json(updatedCottage);
    } catch (err) {
        console.error(err.message);
        res.status(500).send('Server error');
    }
});

// @route   DELETE /api/cottages/:id
// @desc    Delete cottage
// @access  Private
router.delete('/:id', [auth], async (req, res) => {
    try {
        const cottage = await Cottage.findById(req.params.id);
        if (!cottage) {
            return res.status(404).json({ msg: 'Cottage not found' });
        }

        await cottage.remove();
        res.json({ msg: 'Cottage removed' });
    } catch (err) {
        console.error(err.message);
        res.status(500).send('Server error');
    }
});

module.exports = router;
