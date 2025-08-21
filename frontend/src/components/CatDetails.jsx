import './CatDetails.css';

function CatDetails({ cat }) {
  if (!cat) {
    return <p>No cat selected</p>;
  }

  return (
    <div className="cat-details">
    <img 
    src={cat.photo_url} 
    alt={cat.name} 
    style={{ width: '300px', height: 'auto' }} 
    />      
    <div className="details">
        <h2>{cat.name}</h2>
        <p>Date: {cat.date_of_photo}</p>
        <p>Location: {cat.location}</p>
      </div>
    </div>
  );
}

export default CatDetails;
