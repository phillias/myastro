import { useEffect } from 'react';
import { useMap } from 'react-leaflet';
import { GeoSearchControl, OpenStreetMapProvider } from 'leaflet-geosearch';
import 'leaflet/dist/leaflet.css';
import 'leaflet-geosearch/dist/geosearch.css';

interface MapPickerProps {
  onLocationSelect: (lat: number, lon: number) => void;
}

export default function MapPicker({ onLocationSelect }: MapPickerProps) {
  const map = useMap();

  useEffect(() => {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const provider = new (OpenStreetMapProvider as any)();
    
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const searchControl = new (GeoSearchControl as any)({
      provider: provider,
      style: 'bar',
      showMarker: true,
      retainZoomLevel: false,
    });

    map.addControl(searchControl);

    map.on('geosearch/showlocation', (result: any) => {
      onLocationSelect(result.location.y, result.location.x);
    });

    return () => {
      map.removeControl(searchControl);
    };
  }, [map, onLocationSelect]);

  return null;
}
