import React from 'react';

interface NavigationControlsProps {
    onNew: () => void;
    onPrevious: () => void;
    onNext: () => void;
}

const NavigationControls: React.FC<NavigationControlsProps> = ({ onNew, onPrevious, onNext }) => {
    return (
        <div className="flex space-x-3">
            <button
                onClick={onNew}
                className="flex-1 btn-primary"
            >
                NEW
            </button>
            <button
                onClick={onPrevious}
                className="flex-1 btn-secondary"
            >
                PREV
            </button>
            <button
                onClick={onNext}
                className="flex-1 btn-secondary"
            >
                NEXT
            </button>
        </div>
    );
};

export default NavigationControls;

