import React, { useState, useEffect } from 'react';
import './App.css';
import CatDetails from './components/CatDetails';
import UploadModal from './components/UploadModal';

function App() {
  const [currentCat, setCurrentCat] = useState(null);
  const [allCats, setAllCats] = useState([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [hasNext, setHasNext] = useState(false);
  const [hasPrev, setHasPrev] = useState(false);
  const [showUpload, setShowUpload] = useState(false);

  useEffect(() => {
    // initialize page: load a random cat and the list
    init();
  }, []);

  function checkOk(res) {
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    return res.json();
  }

  function applyCatResponse(response) {
    if (!response) return;
    if (response.OK) {
      setCurrentCat(response.Data);
      setHasNext(Boolean(response.Data.has_next));
      setHasPrev(Boolean(response.Data.has_previous));
    }
    else
      setCurrentCat(null);
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
      .then((response) => applyCatResponse(response))
      .catch(() => {/* ignore for now */});
  }

  const handleNext = () => {
    if (!currentCat || !currentCat.created_at) return;
    fetchNext(currentCat.created_at)
      .then((response) => applyCatResponse(response))
      .catch(() => {/* ignore */});
  };

  const handlePrevious = () => {
    if (!currentCat || !currentCat.created_at) return;
    fetchPrevious(currentCat.created_at)
      .then((response) => applyCatResponse(response))
      .catch(() => {/* ignore */});
  };

  return (
    <div className="app">
      <header className="header">
        <h1>CatDex</h1>
        <div className="navigation-buttons">
          <button className='action-button' onClick={() => setShowUpload(true)}>Catch a cat!</button>
        </div>
      </header>

      <main className="main">
        {<CatDetails cat={currentCat} />}

        {showUpload && (
          <UploadModal
            onClose={() => setShowUpload(false)}
            onSuccess={() => {
              setShowUpload(false);
              // refresh random cat
              fetchRandomCat().then((response) => applyCatResponse(response)).catch(() => {});
            }}
          />
        )}

        <div className="navigation-buttons">
          <button className='navigation-button' onClick={handlePrevious} disabled={!hasPrev}>Previous</button>
          <button className='navigation-button' onClick={handleNext} disabled={!hasNext}>Next</button>
        </div>
      </main>
    </div>
  );
}

export default App;
