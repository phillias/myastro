import { useState } from 'react';
import { MapContainer, TileLayer } from 'react-leaflet';
import * as Dialog from '@radix-ui/react-dialog';
import * as Tabs from '@radix-ui/react-tabs';
import { Sparkles, X, MapPin, Calendar, Clock } from 'lucide-react';
import MapPicker from './MapPicker';

interface Remedy {
  mantra: string;
  yantra: string;
  color: string;
}

interface Panchangam {
  tithi: string;
  nakshatra: string;
  yoga: string;
}

interface JyotishData {
  sign: string;
  nakshatra: string;
  lagana_degree: number;
  remedy: Remedy;
  panchanga: Panchangam;
}

export default function App() {
  const [coords, setCoords] = useState<{ lat: number, lon: number } | null>(null);
  const [birthDate, setBirthDate] = useState('');
  const [birthTime, setBirthTime] = useState('');
  const [result, setResult] = useState<JyotishData | null>(null);
  const [open, setOpen] = useState(false);

  const fetchRemedies = async () => {
    if (!coords) return;
    try {
      const url = `/api/calculate?lat=${coords.lat}&lon=${coords.lon}&date=${birthDate}&time=${birthTime}`;
      const response = await fetch(url);
      const data = await response.json();
      setResult(data);
      setOpen(true);
    } catch (err) {
      console.error("Failed to fetch remedies:", err);
    }
  };

  return (
    <div className="min-h-screen bg-slate-900 text-white p-6">
      <header className="flex items-center gap-3 mb-8">
        <Sparkles className="text-purple-400" />
        <h1 className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-purple-400 to-indigo-400">
          Jyotish PWA
        </h1>
      </header>

      <main className="max-w-md mx-auto space-y-6">
        <div className="bg-slate-800 p-6 rounded-2xl border border-white/10 shadow-2xl backdrop-blur-md">
          <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
            <MapPin size={18} /> Enter Birth Details
          </h2>
          
          <div className="space-y-4">
            <div className="h-48 w-full rounded-xl overflow-hidden border border-white/10 relative">
              <MapContainer center={[20, 0]} zoom={2} className="h-full w-full !block">
                <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
                <MapPicker onLocationSelect={(lat, lon) => setCoords({ lat, lon })} />
              </MapContainer>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1">
                <label className="text-xs text-slate-400 flex items-center gap-1"><Calendar size={12}/> Date</label>
                <input type="date" value={birthDate} onChange={(e) => setBirthDate(e.target.value)} className="w-full bg-slate-900 border border-white/10 rounded-lg p-2 text-sm" />
              </div>
              <div className="space-y-1">
                <label className="text-xs text-slate-400 flex items-center gap-1"><Clock size={12}/> Time</label>
                <input type="time" value={birthTime} onChange={(e) => setBirthTime(e.target.value)} className="w-full bg-slate-900 border border-white/10 rounded-lg p-2 text-sm" />
              </div>
            </div>

            <button 
              onClick={fetchRemedies}
              disabled={!coords}
              className="w-full bg-gradient-to-r from-purple-600 to-indigo-600 p-3 rounded-xl font-bold hover:from-purple-500 hover:to-indigo-500 transition shadow-lg disabled:opacity-50"
            >
              Get Jyotish Remedies
            </button>
          </div>
        </div>
      </main>

      <Dialog.Root open={open} onOpenChange={setOpen}>
        <Dialog.Portal>
          <Dialog.Overlay className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50" />
          <Dialog.Content className="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-slate-800 border border-white/10 p-6 rounded-2xl shadow-2xl w-[350px] z-[60]">
            <Dialog.Title className="text-xl font-bold mb-4">Your Cosmic Report</Dialog.Title>
            <Dialog.Description className="hidden">Detailed Jyotish insights</Dialog.Description>
            {result && (
              <Tabs.Root defaultValue="remedy">
                <Tabs.List className="flex border-b border-white/10 mb-4">
                  <Tabs.Trigger value="remedy" className="px-4 py-2 data-[state=active]:border-b-2 data-[state=active]:border-purple-400">Natal</Tabs.Trigger>
                  <Tabs.Trigger value="panchanga" className="px-4 py-2 data-[state=active]:border-b-2 data-[state=active]:border-purple-400">Panchanga</Tabs.Trigger>
                </Tabs.List>
                
                <Tabs.Content value="remedy" className="space-y-2 text-sm">
                  <p className="text-slate-400">Sign: <span className="text-white">{result.sign}</span></p>
                  <div className="mt-4 p-4 bg-white/5 rounded-xl border border-white/5">
                    <p className="font-semibold text-purple-300">Mantra</p>
                    <p className="italic mb-2">{result.remedy.mantra}</p>
                    <p className="font-semibold text-purple-300">Color</p>
                    <p>{result.remedy.color}</p>
                  </div>
                </Tabs.Content>
                
                <Tabs.Content value="panchanga" className="space-y-2 text-sm">
                  <p className="text-slate-400">Tithi: <span className="text-white">{result.panchanga.tithi}</span></p>
                  <p className="text-slate-400">Nakshatra: <span className="text-white">{result.panchanga.nakshatra}</span></p>
                </Tabs.Content>
              </Tabs.Root>
            )}
            <Dialog.Close className="absolute top-4 right-4 text-slate-400 hover:text-white"><X size={20}/></Dialog.Close>
          </Dialog.Content>
        </Dialog.Portal>
      </Dialog.Root>
    </div>
  );
}
