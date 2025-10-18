import axios from 'axios';
import { Conversation, ConversationWithMessages, CreateConversationRequest, Message, SendMessageRequest } from '../types';

const API_BASE_URL = '/api';

const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
        'Content-Type': 'application/json',
    },
});

export const conversationApi = {
    // Get all conversations
    getConversations: async (): Promise<Conversation[]> => {
        const response = await api.get('/conversations');
        return response.data;
    },

    // Get a specific conversation with messages
    getConversation: async (id: string): Promise<ConversationWithMessages> => {
        const response = await api.get(`/conversations/${id}`);
        return response.data;
    },

    // Create a new conversation
    createConversation: async (data: CreateConversationRequest): Promise<Conversation> => {
        const response = await api.post('/conversations', data);
        return response.data.conversation;
    },

    // Send a message to a conversation
    sendMessage: async (conversationId: string, data: SendMessageRequest): Promise<Message> => {
        const response = await api.post(`/conversations/${conversationId}/messages`, data);
        return response.data;
    },

    // Toggle favorite status
    toggleFavorite: async (conversationId: string): Promise<Conversation> => {
        const response = await api.post(`/conversations/${conversationId}/favorite`);
        return response.data;
    },

    // Get next conversation
    getNextConversation: async (currentId: string): Promise<Conversation | null> => {
        try {
            const response = await api.get(`/conversations/${currentId}/next`);
            return response.data;
        } catch (error) {
            return null;
        }
    },

    // Get previous conversation
    getPreviousConversation: async (currentId: string): Promise<Conversation | null> => {
        try {
            const response = await api.get(`/conversations/${currentId}/previous`);
            return response.data;
        } catch (error) {
            return null;
        }
    },
};

