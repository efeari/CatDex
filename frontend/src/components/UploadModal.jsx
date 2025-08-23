import React, { useState, useRef } from 'react';
import './UploadModal.css';

export default function UploadModal({ onClose, onSuccess }) {
  const [file, setFile] = useState(null);
  const [preview, setPreview] = useState(null);
  const [name, setName] = useState('');
  const [date, setDate] = useState(new Date().toISOString().slice(0, 10));
  const [location, setLocation] = useState('');
  const [placeName, setPlaceName] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const inputRef = useRef(null);

  function onFileChosen(f) {
    if (!f) return;
    setFile(f);
    const url = URL.createObjectURL(f);
    setPreview(url);
  }

  function handleDrop(e) {
    e.preventDefault();
    const f = e.dataTransfer?.files?.[0];
    if (f) onFileChosen(f);
  }

  async function reverseGeocode(lat, lon) {
    try {
      const res = await fetch(
        `https://nominatim.openstreetmap.org/reverse?format=jsonv2&lat=${encodeURIComponent(lat)}&lon=${encodeURIComponent(lon)}`,
        { headers: { Accept: 'application/json' } }
      );
      if (!res.ok) throw new Error('Reverse geocode failed');
      const data = await res.json();
      const addr = data.address || {};
      const district = addr.suburb || addr.city_district || addr.county || addr.state_district || '';
      const city = addr.city || addr.town || addr.village || addr.municipality || '';
      const country = addr.country || '';
      const parts = [];
      if (district) parts.push(district);
      if (city && city !== district) parts.push(city);
      if (country) parts.push(country);
      // return the pretty address instead of setting state here; caller will decide where to place it
      return parts.join(', ') || data.display_name || '';
    } catch (err) {
      return '';
    }
  }

  function handleChooseClick() {
    inputRef.current?.click();
  }

  function handleUseLocation() {
    setError(null);
    if (!navigator.geolocation) {
      setError('Geolocation not supported');
      return;
    }
    navigator.geolocation.getCurrentPosition(
      (pos) => {
        const lat = pos.coords.latitude;
        const lon = pos.coords.longitude;
        // swap: put pretty address into the Location input, and show coords (lon,lat) in place-name
        reverseGeocode(lat, lon).then((pretty) => {
          if (pretty) setLocation(pretty);
          else setLocation(`${lon},${lat}`);
          setPlaceName(`${lon},${lat}`);
        }).catch(() => {
          setLocation(`${lon},${lat}`);
          setPlaceName(`${lon},${lat}`);
        });
      },
      () => setError('Could not get location'),
      { enableHighAccuracy: false, timeout: 5000 }
    );
  }

  async function forwardGeocode(query) {
    try {
      setError(null);
      const res = await fetch(
        `https://nominatim.openstreetmap.org/search?format=jsonv2&q=${encodeURIComponent(query)}&limit=1`,
        { headers: { Accept: 'application/json' } }
      );
      if (!res.ok) throw new Error('Forward geocode failed');
      const arr = await res.json();
      if (!arr || arr.length === 0) {
        setError('Location not found');
        return;
      }
      const item = arr[0];
      const lat = item.lat;
      const lon = item.lon;
  // swap: keep original user query (address) in the Location input, show coords (lon,lat) in place-name
  setLocation(query);
  setPlaceName(`${lon},${lat}`);
    } catch (err) {
      setError('Could not find location');
    }
  }

  function handleFindLocation() {
    if (!location || location.trim() === '') {
      setError('Enter a place name to find');
      return;
    }
    // if the user already entered coordinates (lat,lon), skip forward geocoding
    const coordMatch = location.trim().match(/^[-+]?\d+(?:\.\d+)?\s*,\s*[-+]?\d+(?:\.\d+)?$/);
    if (coordMatch) {
      // already coords — clear previous placeName and optionally reverse geocode
      const [lat, lon] = location.split(',').map((s) => s.trim());
      reverseGeocode(lat, lon);
      return;
    }
    forwardGeocode(location.trim());
  }

  async function handleRegister() {
    setError(null);
    if (!file || !name || !date || !location) {
      setError('Please fill all fields and add a photo');
      return;
    }
    setLoading(true);
    try {
      const fd = new FormData();
      fd.append('photo', file);
      fd.append('name', name);
      fd.append('date_of_photo', date);
      fd.append('location', location);
      if (placeName) fd.append('place_name', placeName);

      const res = await fetch('/api/cats', {
        method: 'POST',
        body: fd,
      });
      if (!res.ok) throw new Error(`Upload failed ${res.status}`);
      const body = await res.json();
      setLoading(false);
      onSuccess && onSuccess(body);
      onClose && onClose();
    } catch (err) {
      setLoading(false);
      setError(err.message || 'Upload error');
    }
  }

  return (
    <div className="upload-modal-backdrop" onClick={onClose}>
      <div className="upload-modal" onClick={(e) => e.stopPropagation()} role="dialog">
        <h2>Catch a cat!</h2>

        <div
          className="dropzone"
          onDragOver={(e) => e.preventDefault()}
          onDrop={handleDrop}
          onClick={handleChooseClick}
        >
          {preview ? (
            <img src={preview} alt="preview" className="preview" />
          ) : (
            <div className="drop-instructions">
              <p className="drop-title">Drop photo here</p>
              <p className="drop-sub">or click to choose a file</p>
            </div>
          )}
          <input
            ref={inputRef}
            type="file"
            accept="image/*"
            style={{ display: 'none' }}
            onChange={(e) => onFileChosen(e.target.files?.[0])}
          />
        </div>

        <div className="form-row">
          <label>Name</label>
          <input value={name} onChange={(e) => setName(e.target.value)} />
        </div>

        <div className="form-row">
          <label>Date</label>
          <input type="date" value={date} onChange={(e) => setDate(e.target.value)} />
        </div>

        <div className="form-row">
          <label>Location</label>
          <div className="location-row">
            <input value={location} onChange={(e) => { setLocation(e.target.value); setPlaceName(''); }} placeholder="text or use geolocation" />
            <button type="button" className="small" onClick={handleFindLocation}>Find</button>
            <button type="button" className="small" onClick={handleUseLocation}>Use my location</button>
          </div>
          {placeName && <div className="place-name">{placeName}</div>}
        </div>

        {error && <div className="error">{error}</div>}

        <div className="actions">
          <button className="secondary" onClick={onClose} disabled={loading}>Close</button>
          <button className="primary" onClick={handleRegister} disabled={loading || !file || !name || !date || !location}>
            {loading ? 'Registering…' : 'Register'}
          </button>
        </div>
      </div>
    </div>
  );
}
