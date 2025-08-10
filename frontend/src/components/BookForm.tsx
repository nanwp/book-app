'use client';

import React, { useState, useEffect } from 'react';
import { Book, BookFormData } from '@/types/book';
import { useBooks } from '@/contexts/BookContext';

interface BookFormProps {
  book?: Book;
  onSubmit: (data: BookFormData) => Promise<void>;
  onCancel: () => void;
}

interface FormErrors {
  title?: string;
  author?: string;
  year?: string;
}

const BookForm: React.FC<BookFormProps> = ({ book, onSubmit, onCancel }) => {
  const { loading } = useBooks();
  const [formData, setFormData] = useState<BookFormData>({
    title: book?.title || '',
    author: book?.author || '',
    published_year: book?.published_year || new Date().getFullYear(),
  });
  
  const [errors, setErrors] = useState<FormErrors>({});
  const [touched, setTouched] = useState<Record<string, boolean>>({});

  const validateField = (name: string, value: string | number) => {
    switch (name) {
      case 'title':
        return !value ? 'Book title is required' : '';
      case 'author':
        return !value ? 'Author name is required' : '';
      case 'year':
        const yearNum = Number(value);
        if (!yearNum) return 'Publication year is required';
        if (yearNum < 1000 || yearNum > new Date().getFullYear()) {
          return 'Invalid publication year';
        }
        return '';
      default:
        return '';
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    
    const processedValue = name === 'published_year' ? Number(value) : value;
    
    setFormData(prev => ({
      ...prev,
      [name]: processedValue
    }));

    // Validate field if it has been touched
    if (touched[name]) {
      const error = validateField(name, processedValue);
      setErrors(prev => ({
        ...prev,
        [name]: error
      }));
    }
  };

  const handleBlur = (e: React.FocusEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    const processedValue = name === 'year' ? Number(value) : value;
    
    setTouched(prev => ({ ...prev, [name]: true }));
    
    const error = validateField(name, processedValue);
    setErrors(prev => ({
      ...prev,
      [name]: error
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    // Validate all fields
    const newErrors: FormErrors = {};
    Object.keys(formData).forEach(key => {
      const value = formData[key as keyof BookFormData];
      if (value !== undefined) {
        const error = validateField(key, value);
        if (error) {
          newErrors[key as keyof FormErrors] = error;
        }
      }
    });
    
    setErrors(newErrors);
    setTouched({
      title: true,
      author: true,
      year: true
    });
    
    // If no errors, submit the form
    if (Object.keys(newErrors).length === 0) {
      try {
        await onSubmit(formData);
      } catch (error) {
        console.error('Error submitting form:', error);
      }
    }
  };

  const inputClassName = (fieldName: string) => {
    const baseClass = "w-full px-4 py-3 border rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 transition-all duration-200 bg-gray-50 focus:bg-white text-gray-900 font-medium placeholder:text-gray-500 placeholder:font-normal";
    const hasError = touched[fieldName] && errors[fieldName as keyof FormErrors];
    return `${baseClass} ${hasError ? 'border-red-500 focus:ring-red-500 bg-red-50' : 'border-gray-300 hover:border-gray-400'}`;
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Header */}
      <div className="border-b border-gray-200 pb-4">
        <h3 className="text-lg font-semibold text-gray-900">
          {book ? 'Edit Book Information' : 'Add New Book'}
        </h3>
        <p className="mt-1 text-sm text-gray-600">
          Fill in the book information below. Fields marked with * are required.
        </p>
      </div>

      {/* Basic Information */}
      <div className="space-y-5">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
          {/* Title */}
          <div className="md:col-span-2">
            <label htmlFor="title" className="block text-sm font-semibold text-gray-700 mb-2">
              Book Title *
            </label>
            <input
              type="text"
              id="title"
              name="title"
              value={formData.title}
              onChange={handleInputChange}
              onBlur={handleBlur}
              className={inputClassName('title')}
              placeholder="Enter book title"
            />
            {touched.title && errors.title && (
              <div className="mt-2 flex items-center text-sm text-red-600">
                <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                </svg>
                {errors.title}
              </div>
            )}
          </div>

          {/* Author */}
          <div>
            <label htmlFor="author" className="block text-sm font-semibold text-gray-700 mb-2">
              Author *
            </label>
            <input
              type="text"
              id="author"
              name="author"
              value={formData.author}
              onChange={handleInputChange}
              onBlur={handleBlur}
              className={inputClassName('author')}
              placeholder="Enter author name"
            />
            {touched.author && errors.author && (
              <div className="mt-2 flex items-center text-sm text-red-600">
                <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                </svg>
                {errors.author}
              </div>
            )}
          </div>

          {/* Year */}
          <div>
            <label htmlFor="year" className="block text-sm font-semibold text-gray-700 mb-2">
              Publication Year *
            </label>
            <input
              type="text"
              id="published_year"
              name="published_year"
              value={formData.published_year}
              onChange={handleInputChange}
              onBlur={handleBlur}
              className={inputClassName('year')}
              placeholder="Enter publication year"
            />
            {touched.year && errors.year && (
              <div className="mt-2 flex items-center text-sm text-red-600">
                <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                </svg>
                {errors.year}
              </div>
            )}
          </div>
        </div>

      </div>

      {/* Buttons */}
      <div className="flex justify-end space-x-3 pt-6 border-t border-gray-200">
        <button
          type="button"
          onClick={onCancel}
          className="px-6 py-3 text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-all duration-200 font-medium"
          disabled={loading}
        >
          Cancel
        </button>
        <button
          type="submit"
          className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed font-medium shadow-sm hover:shadow-md"
          disabled={loading}
        >
          {loading ? (
            <div className="flex items-center">
              <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Saving...
            </div>
          ) : (
            book ? 'Update Book' : 'Add Book'
          )}
        </button>
      </div>
    </form>
  );
};

export default BookForm;