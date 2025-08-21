import React, { useState, useEffect } from 'react';
import './App.css';
import CatDetails from './components/CatDetails';

function App() {
  const [currentCat, setCurrentCat] = useState(null);
  const [allCats, setAllCats] = useState([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [hasNext, setHasNext] = useState(false);
  const [hasPrev, setHasPrev] = useState(false);

  useEffect(() => {
    // initialize page: load a random cat and the list
    init();
  }, []);

  function checkOk(res) {
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    return res.json();
  }

  function applyCatResponse(data) {
    if (!data) return;
    setCurrentCat(data);
    setHasNext(Boolean(data.has_next));
    setHasPrev(Boolean(data.has_previous));
  }

  function fetchRandomCat() {
    return fetch('/api/cats/random').then(checkOk);
  }

  function fetchNext(after) {
    return fetch(`/api/cats/next?after=${encodeURIComponent(after)}`).then(checkOk);
  }

  function fetchPrevious(before) {
    return fetch(`/api/cats/previous?before=${encodeURIComponent(before)}`).then(checkOk);
  }

  function init() {
    fetchRandomCat()
      .then((data) => applyCatResponse(data))
      .catch(() => {/* ignore for now */});
  }

  const handleNext = () => {
    if (!currentCat || !currentCat.created_at) return;
    fetchNext(currentCat.created_at)
      .then((data) => applyCatResponse(data))
      .catch(() => {/* ignore */});
  };

  const handlePrevious = () => {
    if (!currentCat || !currentCat.created_at) return;
    fetchPrevious(currentCat.created_at)
      .then((data) => applyCatResponse(data))
      .catch(() => {/* ignore */});
  };

  return (
    <div className="app">
      <header className="header">
        <h1>CatDex</h1>
      </header>

      <main className="main">
        {currentCat && <CatDetails cat={currentCat} />}

        <div className="navigation-buttons">
          <button onClick={handlePrevious} disabled={!hasPrev}>Previous</button>
          <button onClick={handleNext} disabled={!hasNext}>Next</button>
        </div>
      </main>
    </div>
  );
}

export default App;
