import { useEffect, useState } from 'react';
import ChatInterface from './components/ChatInterface';
import LoadingIndicator from './components/LoadingIndicator';
import NavigationControls from './components/NavigationControls';
import SystemPromptEditor from './components/SystemPromptEditor';
import { Conversation, ConversationWithMessages, Message } from './types';
import { conversationApi } from './utils/api';

function App() {
    const [currentConversation, setCurrentConversation] = useState<Conversation | null>(null);
    const [messages, setMessages] = useState<Message[]>([]);
    const [systemPrompt, setSystemPrompt] = useState<string>('');
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    // Load initial conversation on mount
    useEffect(() => {
        loadInitialConversation();
    }, []);

    const loadInitialConversation = async () => {
        try {
            setLoading(true);
            const conversations = await conversationApi.getConversations();
            if (conversations.length > 0) {
                const latestConversation = conversations[conversations.length - 1];
                await loadConversation(latestConversation.id);
            }
        } catch (err) {
            setError('Failed to load conversations');
            console.error('Error loading conversations:', err);
        } finally {
            setLoading(false);
        }
    };

    const loadConversation = async (conversationId: string) => {
        try {
            setLoading(true);
            const data: ConversationWithMessages = await conversationApi.getConversation(conversationId);
            setCurrentConversation(data.conversation);
            setMessages(data.messages);
            setSystemPrompt(data.system_prompt);
            setError(null);
        } catch (err) {
            setError('Failed to load conversation');
            console.error('Error loading conversation:', err);
        } finally {
            setLoading(false);
        }
    };

    const handleNewConversation = async () => {
        try {
            setLoading(true);
            const newConversation = await conversationApi.createConversation({
                title: 'New Conversation',
                system_prompt: systemPrompt || 'You are a helpful AI assistant.',
            });
            setCurrentConversation(newConversation);
            setMessages([]);
            setError(null);
        } catch (err) {
            setError('Failed to create new conversation');
            console.error('Error creating conversation:', err);
        } finally {
            setLoading(false);
        }
    };

    const handleSendMessage = async (content: string) => {
        if (!currentConversation) return;

        try {
            setLoading(true);
            const newMessage = await conversationApi.sendMessage(currentConversation.id, { content });
            setMessages(prev => [...prev, newMessage]);
            setError(null);
        } catch (err) {
            setError('Failed to send message');
            console.error('Error sending message:', err);
        } finally {
            setLoading(false);
        }
    };

    const handleToggleFavorite = async () => {
        if (!currentConversation) return;

        try {
            const updatedConversation = await conversationApi.toggleFavorite(currentConversation.id);
            setCurrentConversation(updatedConversation);
        } catch (err) {
            setError('Failed to toggle favorite');
            console.error('Error toggling favorite:', err);
        }
    };

    const handleNavigate = async (direction: 'previous' | 'next') => {
        if (!currentConversation) return;

        try {
            setLoading(true);
            const conversation = direction === 'next'
                ? await conversationApi.getNextConversation(currentConversation.id)
                : await conversationApi.getPreviousConversation(currentConversation.id);

            if (conversation) {
                await loadConversation(conversation.id);
            }
        } catch (err) {
            setError(`Failed to navigate ${direction}`);
            console.error(`Error navigating ${direction}:`, err);
        } finally {
            setLoading(false);
        }
    };

    const handleSystemPromptChange = (prompt: string) => {
        setSystemPrompt(prompt);
        // Auto-save could be implemented here
    };

    return (
        <div className="flex h-screen bg-gray-50">
            {/* Left Panel - System Prompt Editor */}
            <div className="w-1/2 bg-white border-r border-gray-200 flex flex-col">
                {/* Header */}
                <div className="p-6 border-b border-gray-200 bg-white">
                    <div className="flex items-center justify-between">
                        <h1 className="text-2xl font-bold text-gray-900">Prompt Explorer</h1>
                        <div className="flex items-center space-x-3">
                            <span className="text-sm text-gray-600 font-medium">Conversation ID</span>
                            <span className="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded font-mono">
                                {currentConversation?.id || 'New Conversation'}
                            </span>
                            <button
                                onClick={handleToggleFavorite}
                                className="text-gray-400 hover:text-yellow-500 transition-colors text-xl p-1 rounded hover:bg-yellow-50"
                            >
                                {currentConversation?.is_favorite ? '⭐' : '☆'}
                            </button>
                        </div>
                    </div>
                </div>

                {/* System Prompt Editor */}
                <div className="flex-1 p-6">
                    <SystemPromptEditor
                        value={systemPrompt}
                        onChange={handleSystemPromptChange}
                    />
                </div>

                {/* Navigation Controls */}
                <div className="p-6 border-t border-gray-200 bg-gray-50">
                    <NavigationControls
                        onNew={handleNewConversation}
                        onPrevious={() => handleNavigate('previous')}
                        onNext={() => handleNavigate('next')}
                    />
                </div>
            </div>

            {/* Right Panel - Chat Interface */}
            <div className="w-1/2 bg-gray-100 flex flex-col">
                <ChatInterface
                    messages={messages}
                    onSendMessage={handleSendMessage}
                    loading={loading}
                />
            </div>

            {/* Loading Indicator */}
            {loading && <LoadingIndicator />}

            {/* Error Display */}
            {error && (
                <div className="fixed top-4 right-4 bg-red-500 text-white px-4 py-2 rounded-lg shadow-lg z-50">
                    {error}
                </div>
            )}
        </div>
    );
}

export default App;

