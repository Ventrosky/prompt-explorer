export interface Conversation {
    id: string;
    title: string;
    is_favorite: boolean;
    created_at: string;
}

export interface Message {
    id: string;
    conversation_id: string;
    sender: 'system' | 'user' | 'assistant';
    content: string;
    created_at: string;
}

export interface ConversationWithMessages {
    conversation: Conversation;
    messages: Message[];
    system_prompt: string;
}

export interface CreateConversationRequest {
    title: string;
    system_prompt: string;
}

export interface SendMessageRequest {
    content: string;
}

