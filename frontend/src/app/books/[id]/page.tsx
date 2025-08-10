'use client';

import React, { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import Link from 'next/link';
import { useBooks } from '@/contexts/BookContext';
import { Book } from '@/types/book';
import LoadingSpinner from '@/components/LoadingSpinner';
import ErrorMessage from '@/components/ErrorMessage';
import Modal from '@/components/Modal';
import BookForm from '@/components/BookForm';

export default function BookDetailPage() {
  const params = useParams();
  const router = useRouter();
  const { getBook, updateBook, deleteBook, fetchBooks } = useBooks();
  const [book, setBook] = useState<Book | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);

  const bookId = params.id as string;

  useEffect(() => {
    const loadBook = async () => {
      setLoading(true);
      setError(null);
      try {
        const book = await getBook(bookId);
        if (book) {
          setBook(book);
        } else {
          setError('Buku tidak ditemukan');
        }
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Gagal memuat detail buku');
      } finally {
        setLoading(false);
      }
    };

    if (bookId) {
      loadBook();
    }
  }, [bookId, getBook]);

  const handleEdit = async (bookData: any) => {
    if (book) {
      await updateBook(book.id, bookData);
      setBook({ ...book, ...bookData });
      setIsEditModalOpen(false);
    }
  };

  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

  const handleDelete = async () => {
    await deleteBook(book!.id);
    router.push('/');
    setShowDeleteConfirm(false);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <LoadingSpinner size="lg" message="Memuat detail buku..." />
      </div>
    );
  }

  if (error || !book) {
    return (
      <div className="min-h-screen bg-gray-50">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="mb-6">
            <Link
              href="/"
              className="inline-flex items-center text-blue-600 hover:text-blue-800 transition-colors"
            >
              <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
              Back to Dashboard
            </Link>
          </div>
          <ErrorMessage message={error || 'Buku tidak ditemukan'} />
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="py-6">
            <Link
              href="/"
              className="inline-flex items-center text-blue-600 hover:text-blue-800 transition-colors mb-4"
            >
              <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
              Back to Dashboard
            </Link>
            <h1 className="text-3xl font-bold text-gray-900">Book Detail</h1>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="bg-white rounded-lg shadow-lg overflow-hidden">
          {/* Book Header */}
          <div className="px-6 py-8 border-b border-gray-200">
            <div className="flex justify-between items-start">
              <div className="flex-1">
                <h2 className="text-2xl font-bold text-gray-900 mb-2">{book.title}</h2>
                <p className="text-lg text-gray-600 mb-1">
                  <span className="font-medium">Author:</span> {book.author}
                </p>
                <p className="text-lg text-gray-600">
                  <span className="font-medium">Published Year:</span> {book.published_year}
                </p>
              </div>
              <div className="flex space-x-3">
                <button
                  onClick={() => setIsEditModalOpen(true)}
                  className="px-4 py-2 bg-yellow-100 text-yellow-800 rounded-lg hover:bg-yellow-200 transition-colors font-medium"
                >
                  Edit Book
                </button>
                <button
                  onClick={handleDelete}
                  className="px-4 py-2 bg-red-100 text-red-800 rounded-lg hover:bg-red-200 transition-colors font-medium"
                >
                  Delete Book
                </button>
              </div>
            </div>
          </div>

          {/* Book Details */}
          <div className="px-6 py-8">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
              {/* Left Column */}
            </div>

            <div className="mt-8 pt-8 border-t border-gray-200">
              <div className="bg-gray-50 rounded-lg p-4">
                <h3 className="text-sm font-medium text-gray-500 uppercase tracking-wide mb-2">
                  Additional Information
                </h3>
                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 text-sm">
                  <div>
                    <span className="font-medium text-gray-700">Book ID:</span>
                    <span className="ml-2 text-gray-600 font-mono">{book.id}</span>
                  </div>
                  <div>
                    <span className="font-medium text-gray-700">Year:</span>
                    <span className="ml-2 text-gray-600">{book.published_year}</span>
                  </div>
                  <div>
                    <span className="font-medium text-gray-700">Status:</span>
                    <span className="ml-2 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                      Available
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>

      {/* Edit Modal */}
      <Modal
        isOpen={isEditModalOpen}
        onClose={() => setIsEditModalOpen(false)}
        title="Edit Book"
      >
        <BookForm
          book={book}
          onSubmit={handleEdit}
          onCancel={() => setIsEditModalOpen(false)}
        />
      </Modal>
    </div>
  );
}