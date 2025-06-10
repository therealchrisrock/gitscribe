"use client";

import React, { useState, type ChangeEvent, type FormEvent } from 'react';
import type { FetcherWithComponents } from '@remix-run/react';
import { useFetcher } from '@remix-run/react';

interface CreateMeetingButtonProps {
    className?: string;
    onMeetingCreated?: (meeting: Meeting) => void;
}

interface Meeting {
    id: string;
    title: string;
    type: string;
    status: string;
    start_time: string;
    end_time?: string;
    meeting_url: string;
    recording_path?: string;
    created_at: string;
    updated_at: string;
}

interface CreateMeetingRequest {
    title: string;
    type: 'zoom' | 'google_meet' | 'microsoft_teams' | 'generic';
    meeting_url: string;
}

export function CreateMeetingButton({ className = '', onMeetingCreated }: CreateMeetingButtonProps) {
    const [isOpen, setIsOpen] = useState<boolean>(false);
    const [formData, setFormData] = useState<CreateMeetingRequest>({
        title: '',
        type: 'zoom',
        meeting_url: ''
    });
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setIsLoading(true);
        setError(null);

        try {
            const response = await fetch('/api/meetings', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData),
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Failed to create meeting');
            }

            const meeting: Meeting = await response.json();

            // Reset form and close modal
            setFormData({ title: '', type: 'zoom', meeting_url: '' });
            setIsOpen(false);

            // Notify parent component
            if (onMeetingCreated) {
                onMeetingCreated(meeting);
            }
        } catch (err) {
            setError(err instanceof Error ? err.message : 'An error occurred');
        } finally {
            setIsLoading(false);
        }
    };

    const handleInputChange = (
        e: ChangeEvent<HTMLInputElement | HTMLSelectElement>
    ) => {
        const { name, value } = e.target;

        // Ensure the "type" field maintains its union type.
        const key = name as keyof CreateMeetingRequest;

        setFormData((prev: CreateMeetingRequest) => ({
            ...prev,
            [key]: value,
        } as CreateMeetingRequest));
    };

    return (
        <>
            <button
                onClick={() => setIsOpen(true)}
                className={`
                    bg-blue-600 hover:bg-blue-700 text-white font-medium px-4 py-2 rounded-lg
                    flex items-center gap-2 transition-colors
                    ${className}
                `}
            >
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M12 5V19M5 12H19" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                </svg>
                Create Meeting
            </button>

            {/* Modal */}
            {isOpen && (
                <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
                    <div className="bg-white rounded-lg shadow-xl max-w-md w-full">
                        <div className="p-6">
                            <div className="flex justify-between items-center mb-4">
                                <h2 className="text-xl font-semibold text-gray-900">Create New Meeting</h2>
                                <button
                                    onClick={() => setIsOpen(false)}
                                    className="text-gray-400 hover:text-gray-600"
                                >
                                    <svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                        <path d="M18 6L6 18M6 6L18 18" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                                    </svg>
                                </button>
                            </div>

                            <form onSubmit={handleSubmit} className="space-y-4">
                                <div>
                                    <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-1">
                                        Meeting Title
                                    </label>
                                    <input
                                        type="text"
                                        id="title"
                                        name="title"
                                        value={formData.title}
                                        onChange={handleInputChange}
                                        required
                                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                                        placeholder="Enter meeting title"
                                    />
                                </div>

                                <div>
                                    <label htmlFor="type" className="block text-sm font-medium text-gray-700 mb-1">
                                        Meeting Type
                                    </label>
                                    <select
                                        id="type"
                                        name="type"
                                        value={formData.type}
                                        onChange={handleInputChange}
                                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                                    >
                                        <option value="zoom">Zoom</option>
                                        <option value="google_meet">Google Meet</option>
                                        <option value="microsoft_teams">Microsoft Teams</option>
                                        <option value="generic">Generic</option>
                                    </select>
                                </div>

                                <div>
                                    <label htmlFor="meeting_url" className="block text-sm font-medium text-gray-700 mb-1">
                                        Meeting URL
                                    </label>
                                    <input
                                        type="url"
                                        id="meeting_url"
                                        name="meeting_url"
                                        value={formData.meeting_url}
                                        onChange={handleInputChange}
                                        required={formData.type !== 'generic'}
                                        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                                        placeholder={formData.type === 'generic' ? 'Optional' : 'https://zoom.us/j/123456789'}
                                    />
                                </div>

                                {error && (
                                    <div className="bg-red-50 border border-red-200 rounded-md p-3">
                                        <p className="text-sm text-red-700">{error}</p>
                                    </div>
                                )}

                                <div className="flex gap-3 pt-4">
                                    <button
                                        type="button"
                                        onClick={() => setIsOpen(false)}
                                        className="flex-1 px-4 py-2 text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-md transition-colors"
                                    >
                                        Cancel
                                    </button>
                                    <button
                                        type="submit"
                                        disabled={isLoading}
                                        className="flex-1 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                                    >
                                        {isLoading ? 'Creating...' : 'Create Meeting'}
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>
            )}
        </>
    );
} 