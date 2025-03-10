
from sklearn.metrics import f1_score


class F1():

    def __init__(self):
        self.name = 'f1'
        self.input_type = 'int'
        self.output_type = None
        self.kwargs = None
        

    def parse_data(self, data):
        
        return data



    def compute(self, predictions, references, labels=None, pos_label=1, average="binary", sample_weight=None, **kwargs):
        
        score = f1_score(y_true=references, y_pred=predictions, labels=labels, 
                         pos_label=pos_label, average=average, sample_weight=sample_weight)
        
        score = {
            self.name : float(score) if score.size == 1 else score
            }

        return score