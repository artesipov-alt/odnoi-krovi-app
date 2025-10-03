const express = require('express');
const { register, getProfile, updateProfile, deleteProfile } = require('../controllers/user');
const { createVetClinic, getVetClinics, updateVetClinic, deleteVetClinic } = require('../controllers/vet_clinic');
const { createPet, getPets, updatePet, deletePet } = require('../controllers/pet');
const { createStock, getStocks, updateStock, deleteStock, bookStock } = require('../controllers/blood_stock');
const { createSearch, getSearches, updateSearch, deleteSearch } = require('../controllers/blood_search');
const { planDonation, getDonations, updateDonationStatus, deleteDonation } = require('../controllers/donation');
const { createChat, getChats, deleteChat } = require('../controllers/chat');
const { sendMessage, getMessages, deleteMessage } = require('../controllers/chat_message');

const router = express.Router();

// Users
router.post('/users/register', register);
router.get('/users/profile', getProfile);
router.put('/users/profile', updateProfile);
router.delete('/users/profile', deleteProfile);

// Vet Clinics
router.post('/vet_clinics/create', createVetClinic);
router.get('/vet_clinics', getVetClinics);
router.put('/vet_clinics/:id', updateVetClinic);
router.delete('/vet_clinics/:id', deleteVetClinic);

// Pets
router.post('/pets/create', createPet);
router.get('/pets', getPets);
router.put('/pets/:id', updatePet);
router.delete('/pets/:id', deletePet);

// Blood Stocks
router.post('/blood_stocks/create', createStock);
router.get('/blood_stocks', getStocks);
router.put('/blood_stocks/:id', updateStock);
router.delete('/blood_stocks/:id', deleteStock);
router.post('/blood_stocks/book', bookStock);

// Blood Searches
router.post('/blood_searches/create', createSearch);
router.get('/blood_searches', getSearches);
router.put('/blood_searches/:id', updateSearch);
router.delete('/blood_searches/:id', deleteSearch);

// Donations
router.post('/donations/plan', planDonation);
router.get('/donations', getDonations);
router.put('/donations/:id/status', updateDonationStatus);
router.delete('/donations/:id', deleteDonation);

// Chats
router.post('/chats/create', createChat);
router.get('/chats', getChats);
router.delete('/chats/:id', deleteChat);

// Chat Messages
router.post('/chat_messages/send', sendMessage);
router.get('/chat_messages', getMessages);
router.delete('/chat_messages/:id', deleteMessage);

module.exports = router;