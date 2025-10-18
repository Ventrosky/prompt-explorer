import React from 'react';
import { Message } from '../types';

interface MessageBubbleProps {
    message: Message;
}

const MessageBubble: React.FC<MessageBubbleProps> = ({ message }) => {
    const isUser = message.sender === 'user';
    const isAssistant = message.sender === 'assistant';

    // Don't render system messages in the chat
    if (message.sender === 'system') {
        return null;
    }

    const formatTime = (timestamp: string) => {
        return new Date(timestamp).toLocaleTimeString([], {
            hour: '2-digit',
            minute: '2-digit'
        });
    };

    return (
        <div className={`chat-message flex ${isUser ? 'justify-end' : 'justify-start'} animate-fade-in`}>
            <div className={`message-bubble ${isUser ? 'message-user' : 'message-assistant'}`}>
                <p className="text-sm leading-relaxed">{message.content}</p>
                <p className={`text-xs mt-2 ${isUser ? 'text-blue-200' : 'text-gray-500'}`}>
                    {formatTime(message.created_at)}
                </p>
            </div>
        </div>
    );
};

export default MessageBubble;

