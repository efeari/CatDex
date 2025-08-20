import React, { useState, useEffect } from 'react';
import './App.css';
import CatDetails from './components/CatDetails';

function App() {
  const [randomCat, setRandomCat] = useState(null);
  const [allCats, setAllCats] = useState([]);
  const [currentIndex, setCurrentIndex] = useState(0);

  useEffect(() => {
    // Fetch random cat on page load
    fetch('/api/cats/random')
      .then((response) => response.json())
      .then((data) => setRandomCat(data));

    // Fetch all cats for the list
    // fetch('/api/cats')
    //   .then((response) => response.json())
    //   .then((data) => setAllCats(data));
  }, []);

  const handleNext = () => {
    // if (currentIndex < allCats.length - 1) {
    //   setCurrentIndex(currentIndex + 1);
    //   setRandomCat(allCats[currentIndex + 1]);
    // }
  };

  const handlePrevious = () => {
    // if (currentIndex > 0) {
    //   setCurrentIndex(currentIndex - 1);
    //   setRandomCat(allCats[currentIndex - 1]);
    // }
  };

  return (
    <div className="app">
      <header className="header">
        <h1>CatDex</h1>
      </header>

      <main className="main">
        {randomCat && <CatDetails cat={randomCat} />}

        <div className="navigation-buttons">
          <button onClick={handlePrevious} disabled={currentIndex === 0}>Previous</button>
          <button onClick={handleNext} disabled={currentIndex === allCats.length - 1}>Next</button>
        </div>

        <div className="cat-list">
          <h3>All Cats</h3>
          <ul>
            {allCats.map((cat, index) => (
              <li key={cat.id} onClick={() => {
                setCurrentIndex(index);
                setRandomCat(cat);
              }}>
                <img src={cat.photo_path} alt={cat.name} className="list-cat-image" />
                <span>{cat.name}</span>
              </li>
            ))}
          </ul>
        </div>
      </main>
    </div>
  );
}

export default App;
