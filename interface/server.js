const express = require('express');
const axios = require('axios');
const path = require('path');

const app = express();
const PORT = process.env.PORT || 3000;
const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8000';

app.set('view engine', 'ejs');
app.set('views', path.join(__dirname, 'views'));
app.use(express.urlencoded({ extended: true }));
app.use(express.static('public'));

app.get('/', (req, res) => {
    res.render('index');
});

app.get('/find_order', (req, res) => {
    res.render('index');
});

app.post('/find-order', (req, res) => {
    const orderInfo = req.body.order_info;
    res.redirect(`/order/${orderInfo}`);
});

app.get('/order/:order_id', async (req, res) => {
    try {
        const orderId = req.params.order_id;
        const response = await axios.get(`${API_BASE_URL}/v1/order/get_info/${orderId}`);
        res.render('order', { order: response.data });
    } catch (error) {
        console.error('API Error:', error.message);
        res.render('error', { 
            message: 'Order not found',
            error: error.response?.data || error.message 
        });
    }
});


app.listen(PORT, () => {
    console.log(`Frontend server running on http://localhost:${PORT}`);
});