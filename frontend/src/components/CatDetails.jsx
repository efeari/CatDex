import './CatDetails.css';

function CatDetails({ cat }) {
  if (!cat) {
    return <p>No cat selected</p>;
  }

  return (
  <div className="cat-details">
    <div 
      className="image-container" 
      style={{ width: '400x', height: '400px', overflow: 'hidden' }}
    >
      <img 
        src={cat.photo_url} 
        alt={cat.name} 
        style={{ width: '100%', height: '100%', objectFit: 'cover' }} 
      />
    </div>
    <div className="details">
      <h2>{cat.name}</h2>
      <p>Date: {cat.date_of_photo}</p>
      <p>Location: {cat.location}</p>
    </div>
  </div>

  );
}

export default CatDetails;
