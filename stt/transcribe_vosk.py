import queue
import sys
import json
import sounddevice as sd
from vosk import Model, KaldiRecognizer

# --- Configurações ---
# Caminho para a pasta do modelo Vosk descompactado
MODEL_PATH = "stt/vosk-model-pt-fb-v0.1.1-20220516_2113" # Altere para o caminho correto
# Taxa de amostragem do áudio. 16000 é o padrão para a maioria dos modelos Vosk.
SAMPLERATE = 16000
# --- Fim das Configurações ---

# Fila para armazenar os blocos de áudio do microfone
q = queue.Queue()

def audio_callback(indata, frames, time, status):
    """Esta função é chamada pelo sounddevice para cada bloco de áudio."""
    if status:
        print(status, file=sys.stderr)
    q.put(bytes(indata))

try:
    # Verifica se o modelo existe no caminho especificado
    try:
        model = Model(MODEL_PATH)
    except Exception:
        print(f"Erro: Não foi possível encontrar o modelo em '{MODEL_PATH}'.")
        print("Certifique-se de que você baixou o modelo e o descompactou no caminho correto.")
        sys.exit(1)

    # Cria o objeto de reconhecimento
    recognizer = KaldiRecognizer(model, SAMPLERATE)

    print("Modelo carregado. Pressione Ctrl+C para parar.")
    print("Ouvindo o microfone...")

    # Abre o stream de áudio do microfone
    with sd.RawInputStream(samplerate=SAMPLERATE, blocksize=8000, device=None,
                           dtype='int16', channels=1, callback=audio_callback):

        while True:
            # Pega um bloco de áudio da fila
            data = q.get()
            
            # Alimenta o reconhecedor com o bloco de áudio
            if recognizer.AcceptWaveform(data):
                # Se uma pausa for detectada, imprime o resultado final
                result = json.loads(recognizer.Result())
                if result['text']:
                    print("Transcrição final:", result['text'])
            else:
                # Imprime o resultado parcial enquanto a fala continua
                partial_result = json.loads(recognizer.PartialResult())
                if partial_result['partial']:
                    print("Parcial:", partial_result['partial'])

except KeyboardInterrupt:
    print("\nFinalizando a transcrição.")
except Exception as e:
    print(f"Ocorreu um erro: {type(e).__name__}: {e}")