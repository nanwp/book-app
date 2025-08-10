export interface Book {
  id: string;
  title: string;
  author: string;
  published_year: number;
}

export interface BookFormData {
  title: string;
  author: string;
  published_year: number;
}

export interface BookContextType {
  books: Book[];
  loading: boolean;
  error: string | null;
  addBook: (book: BookFormData) => Promise<void>;
  updateBook: (id: string, book: BookFormData) => Promise<void>;
  deleteBook: (id: string) => Promise<void>;
  fetchBooks: () => Promise<void>;
  getBook: (id: string) => Promise<Book | null>;
}