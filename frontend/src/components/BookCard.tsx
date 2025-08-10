'use client';

import React from 'react';
import Link from 'next/link';
import { Book } from '@/types/book';

interface BookCardProps {
  book: Book;
  onEdit: (book: Book) => void;
  onDelete: (id: string) => void;
}

const BookCard: React.FC<BookCardProps> = ({ book, onEdit, onDelete }) => {
  const handleDelete = () => {
    if (window.confirm(`Apakah Anda yakin ingin menghapus buku "${book.title}"?`)) {
      onDelete(book.id);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow duration-200 p-6 border border-gray-200">
      {/* Book Info */}
      <div className="mb-4">
        <h3 className="text-lg font-semibold text-gray-900 mb-2 line-clamp-2">
          {book.title}
        </h3>
        <p className="text-gray-600 mb-1">
          <span className="font-medium">Author:</span> {book.author}
        </p>
        <p className="text-gray-600 mb-2">
          <span className="font-medium">Year:</span> {book.published_year}
        </p>
      </div>

      {/* Actions */}
      <div className="flex justify-between items-center pt-4 border-t border-gray-100">
        <Link
          href={`/books/${book.id}`}
          className="px-3 py-1 text-sm bg-gray-100 text-gray-700 rounded hover:bg-gray-200 transition-colors"
        >
          Details
        </Link>
        
        <div className="flex space-x-2">
          <button
            onClick={() => onEdit(book)}
            className="px-3 py-1 text-sm bg-blue-100 text-blue-700 rounded hover:bg-blue-200 transition-colors"
          >
            Edit
          </button>
          <button
            onClick={handleDelete}
            className="px-3 py-1 text-sm bg-red-100 text-red-700 rounded hover:bg-red-200 transition-colors"
          >
            Delete
          </button>
        </div>
      </div>
    </div>
  );
};

export default BookCard;