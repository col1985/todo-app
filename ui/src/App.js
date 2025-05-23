import React, { useState, useEffect, useRef } from 'react';
import { PlusCircle, CheckCircle, Circle, Trash2, Edit } from 'lucide-react'; // Icons from lucide-react

// Main App component
const App = () => {
  const [todos, setTodos] = useState([]);
  const [newTodoTitle, setNewTodoTitle] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [editingTodoId, setEditingTodoId] = useState(null);
  const [editingTodoTitle, setEditingTodoTitle] = useState('');

  const inputRef = useRef(null); // Ref for the new todo input field

  // Base URL for the API
  const API_BASE_URL = 'http://localhost:8080/api/v1/todos'; // Ensure this matches your Go API's address

  // Function to fetch all TODOs
  const fetchTodos = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(API_BASE_URL);
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      const data = await response.json();
      setTodos(data);
    } catch (err) {
      console.error("Failed to fetch todos:", err);
      setError("Failed to load todos. Please check if the backend is running.");
    } finally {
      setLoading(false);
    }
  };

  // useEffect to fetch todos on component mount
  useEffect(() => {
    fetchTodos();
  }, []);

  // Function to add a new TODO
  const addTodo = async (e) => {
    e.preventDefault(); // Prevent default form submission
    if (newTodoTitle.trim() === '') {
      setError("Todo title cannot be empty.");
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const response = await fetch(API_BASE_URL, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title: newTodoTitle, completed: false }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const newTodo = await response.json();
      setTodos([...todos, newTodo]);
      setNewTodoTitle(''); // Clear input field
      if (inputRef.current) {
        inputRef.current.focus(); // Focus back on the input
      }
    } catch (err) {
      console.error("Failed to add todo:", err);
      setError("Failed to add todo. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  // Function to toggle TODO completion status
  const toggleTodoCompletion = async (id, currentCompleted) => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`${API_BASE_URL}/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ completed: !currentCompleted }), // Toggle the status
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const updatedTodo = await response.json();
      setTodos(todos.map(todo => (todo.id === id ? updatedTodo : todo)));
    } catch (err) {
      console.error("Failed to update todo:", err);
      setError("Failed to update todo status. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  // Function to delete a TODO
  const deleteTodo = async (id) => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`${API_BASE_URL}/${id}`, {
        method: 'DELETE',
      });

      if (!response.ok && response.status !== 204) { // 204 No Content is expected for successful deletion
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      setTodos(todos.filter(todo => todo.id !== id));
    } catch (err) {
      console.error("Failed to delete todo:", err);
      setError("Failed to delete todo. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  // Function to start editing a TODO
  const startEditing = (todo) => {
    setEditingTodoId(todo.id);
    setEditingTodoTitle(todo.title);
  };

  // Function to save edited TODO title
  const saveEditedTodo = async (id) => {
    if (editingTodoTitle.trim() === '') {
      setError("Todo title cannot be empty.");
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`${API_BASE_URL}/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ title: editingTodoTitle }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const updatedTodo = await response.json();
      setTodos(todos.map(todo => (todo.id === id ? updatedTodo : todo)));
      setEditingTodoId(null); // Exit editing mode
      setEditingTodoTitle('');
    } catch (err) {
      console.error("Failed to save edited todo:", err);
      setError("Failed to update todo title. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  // Function to cancel editing
  const cancelEditing = () => {
    setEditingTodoId(null);
    setEditingTodoTitle('');
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-purple-600 to-indigo-800 flex items-center justify-center p-4 font-inter">
      <div className="bg-white rounded-xl shadow-2xl p-8 w-full max-w-md">
        <h1 className="text-4xl font-extrabold text-center text-gray-800 mb-8">My TODO List</h1>

        {/* Loading and Error Messages */}
        {loading && (
          <div className="text-center text-indigo-600 font-medium mb-4">Loading...</div>
        )}
        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded-lg relative mb-4" role="alert">
            <strong className="font-bold">Error!</strong>
            <span className="block sm:inline ml-2">{error}</span>
          </div>
        )}

        {/* Add New TODO Form */}
        <form onSubmit={addTodo} className="flex items-center gap-2 mb-8">
          <input
            ref={inputRef}
            type="text"
            value={newTodoTitle}
            onChange={(e) => setNewTodoTitle(e.target.value)}
            placeholder="Add a new todo..."
            className="flex-grow p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 transition duration-200"
            aria-label="New todo title"
          />
          <button
            type="submit"
            className="bg-indigo-600 text-white p-3 rounded-lg shadow-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition duration-200 flex items-center justify-center"
            aria-label="Add todo"
          >
            <PlusCircle size={24} />
          </button>
        </form>

        {/* TODO List */}
        {todos.length === 0 && !loading && !error ? (
          <p className="text-center text-gray-500 italic">No todos yet! Add one above.</p>
        ) : (
          <ul className="space-y-4">
            {todos.map((todo) => (
              <li
                key={todo.id}
                className="flex items-center bg-gray-50 p-4 rounded-lg shadow-sm transition duration-200 hover:shadow-md"
              >
                {editingTodoId === todo.id ? (
                  // Editing mode
                  <div className="flex-grow flex items-center gap-2">
                    <input
                      type="text"
                      value={editingTodoTitle}
                      onChange={(e) => setEditingTodoTitle(e.target.value)}
                      onKeyPress={(e) => {
                        if (e.key === 'Enter') saveEditedTodo(todo.id);
                      }}
                      className="flex-grow p-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-1 focus:ring-blue-500"
                    />
                    <button
                      onClick={() => saveEditedTodo(todo.id)}
                      className="bg-green-500 text-white p-2 rounded-lg hover:bg-green-600 transition duration-200"
                      aria-label="Save todo"
                    >
                      Save
                    </button>
                    <button
                      onClick={cancelEditing}
                      className="bg-gray-400 text-white p-2 rounded-lg hover:bg-gray-500 transition duration-200"
                      aria-label="Cancel editing"
                    >
                      Cancel
                    </button>
                  </div>
                ) : (
                  // Display mode
                  <>
                    <button
                      onClick={() => toggleTodoCompletion(todo.id, todo.completed)}
                      className="flex-shrink-0 mr-3 text-gray-400 hover:text-indigo-600 transition duration-200"
                      aria-label={todo.Completed ? "Mark as incomplete" : "Mark as complete"}
                    >
                      {todo.completed ? (
                        <CheckCircle size={28} className="text-green-500" />
                      ) : (
                        <Circle size={28} />
                      )}
                    </button>
                    <span
                      className={`flex-grow text-lg ${todo.completed ? 'line-through text-gray-500' : 'text-gray-800'}`}
                    >
                      {todo.title}
                    </span>
                    <div className="flex-shrink-0 ml-4 flex gap-2">
                      <button
                        onClick={() => startEditing(todo)}
                        className="text-blue-500 hover:text-blue-700 transition duration-200"
                        aria-label="Edit todo"
                      >
                        <Edit size={20} />
                      </button>
                      <button
                        onClick={() => deleteTodo(todo.id)}
                        className="text-red-500 hover:text-red-700 transition duration-200"
                        aria-label="Delete todo"
                      >
                        <Trash2 size={20} />
                      </button>
                    </div>
                  </>
                )}
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
};

export default App;
