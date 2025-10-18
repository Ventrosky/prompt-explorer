import React from 'react';

const LoadingIndicator: React.FC = () => {
    return (
        <div className="fixed top-4 right-4 bg-blue-500 text-white px-4 py-2 rounded-lg shadow-lg z-50">
            <div className="flex items-center space-x-2">
                <div className="loading w-4 h-4 border-2 border-white border-t-transparent rounded-full"></div>
                <span>Processing...</span>
            </div>
        </div>
    );
};

export default LoadingIndicator;

