import React, { useEffect, useState } from 'react';

interface SystemPromptEditorProps {
    value: string;
    onChange: (value: string) => void;
}

const SystemPromptEditor: React.FC<SystemPromptEditorProps> = ({ value, onChange }) => {
    const [localValue, setLocalValue] = useState(value);

    useEffect(() => {
        setLocalValue(value);
    }, [value]);

    const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        const newValue = e.target.value;
        setLocalValue(newValue);
        onChange(newValue);
    };

    const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
        // Auto-resize functionality
        const textarea = e.target as HTMLTextAreaElement;
        textarea.style.height = 'auto';
        textarea.style.height = textarea.scrollHeight + 'px';
    };

    return (
        <div className="h-full flex flex-col">
            <label className="block text-sm font-medium text-gray-700 mb-3">
                System Prompt
            </label>
            <textarea
                value={localValue}
                onChange={handleChange}
                onKeyDown={handleKeyDown}
                placeholder="Enter your system prompt here..."
                className="flex-1 w-full p-4 border-2 border-gray-300 rounded-lg resize-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all font-mono text-sm bg-gray-50 focus:bg-white"
                style={{ minHeight: '200px' }}
            />
        </div>
    );
};

export default SystemPromptEditor;

