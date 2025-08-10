'use client';

import React, { createContext, useContext, useState, useCallback, ReactNode } from 'react';
import { Book, BookFormData, BookContextType } from '@/types/book';

const BookContext = createContext<BookContextType | undefined>(undefined);

const baseUrl = process.env.NEXT_PUBLIC_API_URL;

export const useBooks = () => {
  const context = useContext(BookContext);
  if (context === undefined) {
    throw new Error('useBooks must be used within a BookProvider');
  }
  return context;
};

interface BookProviderProps {
  children: ReactNode;
}

export const BookProvider: React.FC<BookProviderProps> = ({ children }) => {
  const [books, setBooks] = useState<Book[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchBooks = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`${baseUrl}/api/v1/books`);
      if (!response.ok) {
        throw new Error('Failed to fetch books');
      }
      const data = await response.json();
      setBooks(data.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
      setBooks([
        {
          id: '1',
          title: 'The Great Gatsby',
          author: 'F. Scott Fitzgerald',
          published_year: 1925,
        },
        {
          id: '2',
          title: 'To Kill a Mockingbird',
          author: 'Harper Lee',
          published_year: 1960,
        }
      ]);
    } finally {
      setLoading(false);
    }
  }, []);

  const addBook = useCallback(async (bookData: BookFormData) => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`${baseUrl}/api/v1/books`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(bookData),
      });
      
      if (!response.ok) {
        throw new Error('Failed to add book');
      }
      
      const newBook = await response.json();
      setBooks(prev => [newBook.data, ...prev]);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to add book');
    } finally {
      setLoading(false);
    }
  }, []);

  const updateBook = useCallback(async (id: string, bookData: BookFormData) => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`${baseUrl}/api/v1/books/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(bookData),
      });
      
      if (!response.ok) {
        throw new Error('Failed to update book');
      }
      
      const updatedBook = await response.json();
      setBooks(prev => prev.map(book => book.id === id ? updatedBook.data : book));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update book');
    } finally {
      setLoading(false);
    }
  }, []);

  const deleteBook = useCallback(async (id: string) => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`${baseUrl}/api/v1/books/${id}`, {
        method: 'DELETE',
      });
      
      if (!response.ok) {
        throw new Error('Failed to delete book');
      }
      
      setBooks(prev => prev.filter(book => book.id !== id));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete book');
      // Fallback untuk demo
    } finally {
      setLoading(false);
    }
  }, []);

  const getBook = useCallback(async (id: string) => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`${baseUrl}/api/v1/books/${id}`);
      if (!response.ok) {
        throw new Error('Failed to fetch book');
      }
      const data = await response.json();
      return data.data;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch book');
      return null;
    } finally {
      setLoading(false);
    }
  }, [books]);

  const value: BookContextType = {
    books,
    loading,
    error,
    addBook,
    updateBook,
    deleteBook,
    getBook,
    fetchBooks,
  };

  return (
    <BookContext.Provider value={value}>
      {children}
    </BookContext.Provider>
  );
};